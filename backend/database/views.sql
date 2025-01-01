CREATE VIEW contender_score AS (
  SELECT score.* FROM contender LEFT JOIN score ON score.contender_id = contender.id
);