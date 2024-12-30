-- name: GetContender :one
SELECT sqlc.embed(contender), sqlc.embed(score) FROM contender
JOIN score ON score.contender_id = id
WHERE id = ?;

-- name: GetContenderByCode :one
SELECT sqlc.embed(contender), sqlc.embed(score) FROM contender
JOIN score ON score.contender_id = id
WHERE registration_code = ?;

-- name: GetContendersByCompClass :many
SELECT sqlc.embed(contender), sqlc.embed(score) FROM contender
JOIN score ON score.contender_id = id
WHERE class_id = ?;

-- name: GetContendersByContest :many
SELECT sqlc.embed(contender), sqlc.embed(score) FROM contender
JOIN score ON score.contender_id = id
WHERE contest_id = ?;

-- name: DeleteContender :exec
DELETE FROM contender
WHERE id = ?;

-- name: CountContenders :one
SELECT COUNT(*)
FROM contender
WHERE contest_id = ?;

-- name: UpsertContender :execlastid
INSERT INTO 
	contender (id, organizer_id, contest_id, registration_code, name, club, class_id, entered, disqualified, withdrawn_from_finals)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id),
    registration_code = VALUES(registration_code),
    name = VALUES(name),
    club = VALUES(club),
    class_id = VALUES(class_id),
    entered = VALUES(entered),
    disqualified = VALUES(disqualified),
    withdrawn_from_finals = VALUES(withdrawn_from_finals);

-- name: UpsertScore :exec
INSERT INTO
    score (contender_id, timestamp, score, placement, finalist, rank_order)
VALUES
    (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    timestamp = VALUES(timestamp),
    score = VALUES(score),
    placement = VALUES(placement),
    finalist = VALUES(finalist),
    rank_order = VALUES(rank_order);

-- name: GetCompClass :one
SELECT sqlc.embed(comp_class)
FROM comp_class
WHERE id = ?;

-- name: GetCompClassesByContest :many
SELECT sqlc.embed(comp_class)
FROM comp_class
WHERE contest_id = ?;

-- name: GetContest :one
SELECT sqlc.embed(contest), MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
FROM contest
LEFT JOIN comp_class cc ON cc.contest_id = contest.id
WHERE contest.id = ?;

-- name: GetContestsCurrentlyRunningOrByStartTime :many
SELECT sqlc.embed(contest), MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
FROM contest
JOIN comp_class cc ON cc.contest_id = contest.id
GROUP BY contest.id
HAVING
    NOW() BETWEEN MIN(cc.time_begin) AND MAX(cc.time_end)
	OR MIN(cc.time_begin) BETWEEN sqlc.arg(earliestStartTime) AND sqlc.arg(latestStartTime);

-- name: GetProblem :one
SELECT sqlc.embed(problem)
FROM problem
WHERE id = ?;

-- name: GetProblemsByContest :many
SELECT sqlc.embed(problem)
FROM problem
WHERE contest_id = ?;

-- name: GetTick :one
SELECT sqlc.embed(tick)
FROM tick
WHERE id = ?;

-- name: GetTicksByContender :many
SELECT sqlc.embed(tick)
FROM tick
WHERE contender_id = ?;

-- name: GetTicksByContest :many
SELECT sqlc.embed(tick)
FROM tick
WHERE contest_id = ?;

-- name: DeleteTick :exec
DELETE
FROM tick
WHERE id = ?;

-- name: InsertTick :execlastid
INSERT INTO
    tick (organizer_id, contest_id, contender_id, problem_id, flash, timestamp)
VALUES
    (?, ?, ?, ?, ?, ?);