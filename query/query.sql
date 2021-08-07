DECLARE timezone STRING;
DECLARE date_from DATE;
DECLARE date_to DATE;

SET timezone = 'Asia/Tokyo';
SET date_from = DATE(TIMESTAMP('{{.ReportingDateFrom}}'), timezone);
SET date_to = DATE(TIMESTAMP('{{.ReportingDateTo}}'), timezone);

WITH
  this_month AS(
  SELECT
    service.description AS service,
    cost AS monthly,
    CASE
      WHEN DATE(usage_end_time, timezone) = date_to THEN cost
    ELSE
      0
    END
    AS yesterday
  FROM
    `{{.TableName}}`
  WHERE
    DATE(_PARTITIONTIME, timezone) BETWEEN date_from AND date_to
    AND DATE(usage_end_time, timezone) BETWEEN date_from AND date_to),
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