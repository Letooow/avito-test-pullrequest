-- name: SaveUser :exec
INSERT INTO users (userid, username, isactive)
VALUES ($1, $2, $3);

-- name: UpdateUser :exec
UPDATE users SET username = $1, isactive = $2 WHERE userid = $3;

-- name: GetUserByID :one
SELECT * FROM users WHERE userid = $1;

-- name: GetUsersByTeamName :many
SELECT u.userid,
       u.username,
       u.isactive
FROM users u
         JOIN users_team ut ON ut.userid = u.userid
WHERE ut.teamname = $1;

-- name: GetUsersTeams :many
SELECT t.teamname
FROM teams t
         JOIN users_team ut ON ut.teamname = t.teamname
WHERE ut.userid = $1;

-- name: SaveUserTeam :exec
INSERT INTO users_team (teamname, userid) VALUES ($1, $2);

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateTeam :exec
INSERT INTO teams (teamname) VALUES ($1);

-- name: GetTeamByName :one
SELECT * FROM teams WHERE teamname = $1;

-- name: GetTeams :many
SELECT teamname FROM teams;

-- name: CreatePullRequest :exec
INSERT INTO pull_requests (pullrequestid, name, status) VALUES ($1, $2, $3);

-- name: UpdatePullRequestStatus :exec
UPDATE pull_requests SET status = $1, mergedat = $2 WHERE pullrequestid = $3;

-- name: AssignUserPullRequest :exec
INSERT INTO users_pull_requests (pullrequestid, userid, role) VALUES ($1, $2, $3);

-- name: GetPullRequestByID :one
SELECT * FROM pull_requests WHERE pullrequestid = $1;

-- name: GetPullRequests :many
SELECT * FROM pull_requests;

-- name: GetUsersAssignedPullRequest :many
SELECT pull_requests.pullrequestid, role FROM pull_requests, users_pull_requests WHERE pull_requests.pullrequestid = (SELECT pullrequestid FROM users_pull_requests WHERE users_pull_requests.userid = $1);

-- name: GetListOfUsersByPullRequestID :many
SELECT userid, role FROM users_pull_requests WHERE pullrequestid = $1;

-- name: DeletePullRequestAssignOfUser :exec
DELETE FROM users_pull_requests WHERE pullrequestid = $1 AND userid = $2;

