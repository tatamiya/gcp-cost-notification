WITH
  this_month AS(
  SELECT
    service.description AS service,
    cost AS monthly,
    CASE
      WHEN DATE(usage_end_time) = DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY) THEN cost
    ELSE
    0
  END
    AS yesterday
  FROM
    {{.TableName}}
  WHERE
    DATE(_PARTITIONTIME) >= DATE_TRUNC(CURRENT_DATE(), MONTH)
    AND DATE(usage_end_time) >= DATE_TRUNC(CURRENT_DATE(), MONTH) ),
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