CREATE TABLE users
(
    UserID   TEXT UNIQUE NOT NULL PRIMARY KEY,
    Username VARCHAR(128) NOT NULL,
    IsActive BOOL        NOT NULL DEFAULT TRUE
);


CREATE TABLE teams
(
    TeamName TEXT UNIQUE NOT NULL PRIMARY KEY
);

CREATE TABLE pull_requests
(
    PullRequestID TEXT UNIQUE NOT NULL PRIMARY KEY,
    Name          TEXT,
    Status        VARCHAR(32) NOT NULL,
    CreatedAt     TIMESTAMP NOT NULL DEFAULT now(),
    MergedAt      TIMESTAMP
);

CREATE TABLE users_pull_requests
(
    PullRequestID TEXT NOT NULL REFERENCES pull_requests(PullRequestID),
    UserID TEXT NOT NULL REFERENCES users(UserID),
    Role VARCHAR(64) NOT NULL,
    PRIMARY KEY (PullRequestID, UserID, Role)
);

CREATE INDEX idx_upr_user_id ON users_pull_requests (UserID);

CREATE TABLE users_team
(
    TeamName TEXT NOT NULL REFERENCES teams(TeamName),
    UserID TEXT NOT NULL REFERENCES users(UserID),
    PRIMARY KEY (TeamName, UserID)
);

CREATE INDEX idx_ut_user_id ON users_team (UserID);

