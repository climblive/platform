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
	contender (id, organizer_id, contest_id, registration_code, name, class_id, entered, disqualified, withdrawn_from_finals)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id),
    registration_code = VALUES(registration_code),
    name = VALUES(name),
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

-- name: GetAllContests :many
SELECT sqlc.embed(contest), MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end
FROM contest
LEFT JOIN comp_class cc ON cc.contest_id = contest.id
GROUP BY contest.id;

-- name: UpsertContest :execlastid
INSERT INTO 
	contest (id, organizer_id, archived, series_id, name, description, location, qualifying_problems, finalists, rules, grace_period, created)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    archived = VALUES(archived),
    series_id = VALUES(series_id),
    name = VALUES(name),
    description = VALUES(description),
    location = VALUES(location),
    qualifying_problems = VALUES(qualifying_problems),
    finalists = VALUES(finalists),
    rules = VALUES(rules),
    grace_period = VALUES(grace_period),
    created = VALUES(created);

-- name: DeleteContest :exec
DELETE FROM contest
WHERE id = ?;

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
    WHERE archived = FALSE
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
	problem (id, organizer_id, contest_id, number, hold_color_primary, hold_color_secondary, zone_1_enabled, zone_2_enabled, description, points_zone_1, points_zone_2, points_top, flash_bonus)
VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id),
    number = VALUES(number),
    hold_color_primary = VALUES(hold_color_primary),
    hold_color_secondary = VALUES(hold_color_secondary),
    zone_1_enabled = VALUES(zone_1_enabled),
    zone_2_enabled = VALUES(zone_2_enabled),
    description = VALUES(description),
    points_zone_1 = VALUES(points_zone_1),
    points_zone_2 = VALUES(points_zone_2),
    points_top = VALUES(points_top),
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
    tick (organizer_id, contest_id, contender_id, problem_id, timestamp, top, attempts_top, zone_1, attempts_zone_1, zone_2, attempts_zone_2)
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpsertOrganizer :execlastid
INSERT INTO
    organizer (id, name)
VALUES
    (?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name);

-- name: UpsertUser :execlastid
INSERT INTO
    user (id, username, admin)
VALUES
    (?, ?, ?)
ON DUPLICATE KEY UPDATE
    username = VALUES(username),
    admin = VALUES(admin);

-- name: GetUserByUsername :many
SELECT sqlc.embed(user), sqlc.embed(organizer)
FROM user
LEFT JOIN user_organizer uo ON uo.user_id = user.id
LEFT JOIN organizer ON organizer.id = uo.organizer_id
WHERE username = ?;

-- name: GetUsersByOrganizer :many
SELECT sqlc.embed(user)
FROM user
LEFT JOIN user_organizer uo ON uo.user_id = user.id
WHERE uo.organizer_id = ?;

-- name: AddUserToOrganizer :exec
INSERT INTO
    user_organizer (user_id, organizer_id)
VALUES
    (?, ?);

-- name: GetOrganizer :one
SELECT *
FROM organizer
WHERE id = ?;

-- name: GetAllOrganizers :many
SELECT *
FROM organizer;

-- name: GetRaffle :one
SELECT sqlc.embed(raffle)
FROM raffle
WHERE id = ?;

-- name: GetRafflesByContest :many
SELECT sqlc.embed(raffle)
FROM raffle
WHERE contest_id = ?;

-- name: UpsertRaffle :execlastid
INSERT INTO
    raffle (id, organizer_id, contest_id)
VALUES
    (?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    contest_id = VALUES(contest_id);

-- name: GetRaffleWinners :many
SELECT sqlc.embed(raffle_winner), contender.name
FROM raffle_winner
JOIN contender ON contender.id = raffle_winner.contender_id
WHERE raffle_id = ?;

-- name: UpsertRaffleWinner :execlastid
INSERT INTO
    raffle_winner (id, organizer_id, raffle_id, contender_id, timestamp)
VALUES
    (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    organizer_id = VALUES(organizer_id),
    raffle_id = VALUES(raffle_id),
    contender_id = VALUES(contender_id),
    timestamp = VALUES(timestamp); 

-- name: GetOrganizerInvitesByOrganizer :many
SELECT sqlc.embed(organizer_invite), organizer.name
FROM organizer_invite
JOIN organizer ON organizer.id = organizer_invite.organizer_id
WHERE organizer_id = ?;

-- name: GetOrganizerInvite :one
SELECT sqlc.embed(organizer_invite), organizer.name
FROM organizer_invite
JOIN organizer ON organizer.id = organizer_invite.organizer_id
WHERE organizer_invite.id = ?;

-- name: InsertOrganizerInvite :exec
INSERT INTO
    organizer_invite (id, organizer_id, expires_at)
VALUES
    (?, ?, ?);

-- name: DeleteOrganizerInvite :exec
DELETE FROM organizer_invite
WHERE id = ?;
