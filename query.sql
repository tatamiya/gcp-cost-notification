DECLARE timezone STRING;
DECLARE execution_date DATE;
DECLARE one_day_before DATE;

SET timezone = 'Asia/Tokyo';
SET execution_date = DATE(TIMESTAMP({{.ExecutionTimestamp}}), timezone);
SET one_day_before = DATE_SUB(execution_date, INTERVAL 1 DAY);

WITH
  this_month AS(
  SELECT
    service.description AS service,
    cost AS monthly,
    CASE
      WHEN DATE(usage_end_time, timezone) = one_day_before THEN cost
    ELSE
    0
  END
    AS yesterday
  FROM
    {{.TableName}}
  WHERE
    DATE(_PARTITIONTIME, timezone) BETWEEN DATE_TRUNC(one_day_before, MONTH) AND one_day_before
    AND DATE(usage_end_time, timezone) BETWEEN DATE_TRUNC(one_day_before, MONTH) AND one_day_before),
  details AS (
  SELECT
    service,
    ROUND(SUM(monthly),2) AS monthly,
    ROUND(SUM(yesterday),2) AS yesterday
  FROM
    this_month
  GROUP BY
    service
  HAVING
    monthly > 0 )
SELECT
  'Total' AS service,
  ROUND(SUM(monthly),2) AS monthly,
  ROUND(SUM(yesterday),2) AS yesterday
FROM
  this_month
UNION ALL
SELECT
  service,
  monthly,
  yesterday
FROM
  details
ORDER BY
  monthly DESC