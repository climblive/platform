-- name: GetContender :one
SELECT sqlc.embed(contender), score.*
FROM contender
LEFT JOIN score ON score.contender_id = id
WHERE id = ?;

-- name: GetContenderByCode :one
SELECT sqlc.embed(contender), score.*
FROM contender
LEFT JOIN score ON score.contender_id = id
WHERE registration_code = ?;

-- name: GetContendersByCompClass :many
SELECT sqlc.embed(contender), score.*
FROM contender
LEFT JOIN score ON score.contender_id = id
WHERE class_id = ?;

-- name: GetContendersByContest :many
SELECT sqlc.embed(contender), score.*
FROM contender
LEFT JOIN score ON score.contender_id = id
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

-- name: DeleteCompClass :exec
DELETE FROM comp_class
WHERE id = ?;

-- name: UpsertCompClass :execlastid
INSERT INTO 
	comp_class (id, organizer_id, contest_id, name, description, color, time_begin, time_end)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id),
    name = VALUES(name),
    description = VALUES(description),
    color = VALUES(color),
    time_begin = VALUES(time_begin),
    time_end = VALUES(time_end);

-- name: GetContest :one
SELECT sqlc.embed(contest), MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
FROM contest
LEFT JOIN comp_class cc ON cc.contest_id = contest.id
WHERE contest.id = ?
GROUP BY contest.id;

-- name: UpsertContest :execlastid
INSERT INTO 
	contest (id, organizer_id, series_id, name, description, location, final_enabled, qualifying_problems, finalists, rules, grace_period)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    series_id = VALUES(series_id),
    name = VALUES(name),
    description = VALUES(description),
    location = VALUES(location),
    final_enabled = VALUES(final_enabled),
    qualifying_problems = VALUES(qualifying_problems),
    finalists = VALUES(finalists),
    rules = VALUES(rules),
    grace_period = VALUES(grace_period);

-- name: GetContestsByOrganizer :many
SELECT sqlc.embed(contest), MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
FROM contest
LEFT JOIN comp_class cc ON cc.contest_id = contest.id
WHERE contest.organizer_id = ?
GROUP BY contest.id;

-- name: GetContestsCurrentlyRunningOrByStartTime :many
SELECT
	*
FROM (
    SELECT contest.*, MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
    FROM contest
    JOIN comp_class cc ON cc.contest_id = contest.id
    GROUP BY contest.id) AS sub
WHERE
    NOW() BETWEEN sub.time_begin AND DATE_ADD(sub.time_end, INTERVAL (sub.grace_period + 15) MINUTE)
	OR sub.time_begin BETWEEN sqlc.arg(earliest_start_time) AND sqlc.arg(latest_start_time);

-- name: GetProblem :one
SELECT sqlc.embed(problem)
FROM problem
WHERE id = ?;

-- name: GetProblemByNumber :one
SELECT sqlc.embed(problem)
FROM problem
WHERE contest_id = ? AND number = ?;

-- name: GetProblemsByContest :many
SELECT sqlc.embed(problem)
FROM problem
WHERE contest_id = ?;

-- name: DeleteProblem :exec
DELETE FROM problem
WHERE id = ?;

-- name: UpsertProblem :execlastid
INSERT INTO 
	problem (id, organizer_id, contest_id, number, hold_color_primary, hold_color_secondary, name, description, points, flash_bonus)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id),
    number = VALUES(number),
    hold_color_primary = VALUES(hold_color_primary),
    hold_color_secondary = VALUES(hold_color_secondary),
    name = VALUES(name),
    description = VALUES(description),
    points = VALUES(points),
    flash_bonus = VALUES(flash_bonus);

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

-- name: GetTicksByProblem :many
SELECT sqlc.embed(tick)
FROM tick
WHERE problem_id = ?;

-- name: DeleteTick :exec
DELETE
FROM tick
WHERE id = ?;

-- name: InsertTick :execlastid
INSERT INTO
    tick (organizer_id, contest_id, contender_id, problem_id, flash, timestamp)
VALUES
    (?, ?, ?, ?, ?, ?);

-- name: UpsertOrganizer :execlastid
INSERT INTO
    organizer (id, name, homepage)
VALUES
    (?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    homepage = VALUES(homepage);

-- name: UpsertUser :execlastid
INSERT INTO
    user (id, name, username, admin)
VALUES
    (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    username = VALUES(username),
    admin = VALUES(admin);

-- name: GetUserByUsername :many
SELECT sqlc.embed(user), organizer.id AS organizer_id
FROM user
LEFT JOIN user_organizer uo ON uo.user_id = user.id
LEFT JOIN organizer ON organizer.id = uo.organizer_id
WHERE username = ?;

-- name: AddUserToOrganizer :exec
INSERT INTO
    user_organizer (user_id, organizer_id)
VALUES
    (?, ?);

-- name: GetOrganizer :one
SELECT *
FROM organizer
WHERE id = ?;
