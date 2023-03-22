-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO username;
CREATE TYPE bracket_type AS ENUM ('SINGLE_ELIMINATION', 'DOUBLE_ELIMINATION', 'ROUND_ROBIN');
CREATE TYPE bracket_status AS ENUM ('pending', 'live', 'finished');

create table if not exists users_email(
    user_id bigint not null references user_accounts(id) on delete cascade,
    email varchar(254) not null unique primary key,
    is_verified boolean default false not null,
    url_verification varchar(20),
    change_at timestamp
);
create table if not exists user_accounts(
    id bigint primary key generated always as identity,
    user_login varchar(32) not null unique
--     email varchar(254) not null references users_email(email)
);

create table if not exists user_tokens(
    user_id bigint references user_accounts(id) not null,
    access_token varchar(20),
    refresh_token varchar(20),
    expire_at timestamp not null,
    constraint user_key primary key (user_id)
);
create table if not exists tournaments(
    id uuid not null primary key,
    sport_name varchar(32) not null,
    title varchar(32) NOT NULL,
    start_at timestamp NOT NULL,
    end_at timestamp NOT NULL check ( end_at > start_at ),
    description json,
    created_by_user bigint references user_accounts(id) on update cascade,
    brackets_limit int default 1 not null check ( brackets_limit > 0 ),
    is_active boolean default true not null
);

create table if not exists brackets(
    id uuid not null primary key,
    type_of bracket_type NOT NULL,
    max_teams int check (max_teams > 0 and max_teams % 2 = 0),
    max_team_participants int not null default 1 check ( max_team_participants > 0 ),
    tournament_id uuid not null references tournaments(id) on delete cascade,
    playoff_rounds int default 1 NOT NULL check ( playoff_rounds > 0 ),
    final_rounds int default 1 NOT NULL check ( playoff_rounds > 0 ),
    grand_final_rounds int default 1 NOT NULL check ( playoff_rounds > 0 ),
    status bracket_status default 'pending' not null
);
create table if not exists teams(
    id bigint primary key generated always as identity,
    team_alias varchar(20) not null,
    bracket_id uuid references brackets(id) on delete cascade,
    unique (team_alias, bracket_id)
);
create table if not exists matches(
    id bigint primary key generated always as identity,
    bracket_id uuid references brackets(id) on delete cascade,
    round int not null, --check ( round % 2 = 0 ),
    first_team bigint references teams(id) on DELETE set null on UPDATE cascade,
    second_team bigint references teams(id) on DELETE set null on UPDATE cascade,
    first_team_score int,
    second_team_score int,
    start_on timestamp,
    winner bigint references teams(id) on update CASCADE on DELETE set null,
--     prev_firts_team_match bigint references matches(id)
    unique (bracket_id, round, first_team, second_team)
);
create table if not exists participants(
    id bigint primary key generated always as identity,
    user_alias varchar(20) not null,
    team_id bigint references teams(id) on delete cascade,
    contact varchar(100),
    unique (team_id, user_alias)
);
