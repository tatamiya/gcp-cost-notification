DECLARE execution_date DATE DEFAULT CURRENT_DATE();
SET execution_date = DATE(TIMESTAMP {{.ExecutionTimestamp}});

WITH
  this_month AS(
  SELECT
    service.description AS service,
    cost AS monthly,
    CASE
      WHEN DATE(usage_end_time) = DATE_SUB(execution_date, INTERVAL 1 DAY) THEN cost
    ELSE
    0
  END
    AS yesterday
  FROM
    {{.TableName}}
  WHERE
    DATE(_PARTITIONTIME) >= DATE_TRUNC(execution_date, MONTH)
    AND DATE(usage_end_time) >= DATE_TRUNC(execution_date, MONTH) ),
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