drop table if exists test;
-- CREATE TYPE tournament_system AS ENUM ('elimination', 'swiss', 'robin');
-- CREATE TYPE bracket_format AS ENUM ('GROUP_STAGE');
-- single, double elimination; Round Robin
CREATE TYPE bracket_type AS ENUM ('SINGLE_ELIMINATION', 'DOUBLE_ELIMINATION', 'ROUND_ROBIN');
CREATE TYPE bracket_status AS ENUM ('pending', 'live', 'finished');

create table if not exists users(
    id bigserial not null generated always as identity primary key,
    login varchar(32) not null unique,
    first_name varchar(32),
    last_name varchar(32),
    email varchar(254),
    secure_key varchar(32)
--     token_fk
);
create table if not exists users_email(
    email varchar(254) not null check (  ),
);
create table if not exists tokens(
    refresh_token varchar(0)
);
create table if not exists matches(
    id bigserial not null unique generated always as identity primary key,
    bracket_id uuid references brackets(id) on delete cascade,
    round int not null check ( round % 2 == 0 ),
    first_team bigserial references teams(id) on DELETE set null,
    second_team bigserial references teams(id) on DELETE set null,
    first_team_score int,
    second_team_score int,
    start_on timestamp,
    winner bigserial references teams(id) on update CASCADE on DELETE set null
);
create table if not exists teams(
    id bigserial not null unique generated always as identity primary key,
    team_alias varchar(20),
    bracket_id uuid references brackets(id) on delete cascade
);
create table if not exists participants(
    id bigserial not null unique generated always as identity primary key,
    user_alias varchar(20) not null,
    team_id bigserial references teams(id) on delete cascade,
    contact text
);
create table if not exists brackets(
    id uuid unique not null primary key,
    type_of bracket_type NOT NULL,
    teams_limit int check (teams_limit > 0 and teams_limit % 2 == 0),
    participants_limit int not null default 1 check ( participants_limit > 0 ),
    tournament_id uuid references tournaments(id) on delete cascade,
    playoff_rounds int default 1 NOT NULL check ( playoff_rounds > 0 ),
    final_rounds int default 1 not null check ( playoff_rounds > 0 ),
    grand_final_rounds int default 1 not null check ( playoff_rounds > 0 ),
    status bracket_status default 'pending' not null
);
create table if not exists tournaments(
    id uuid unique not null primary key,
    sport_name varchar(32) not null,
    title varchar(32) NOT NULL,
    start_at timestamp default NOT NULL,
    end_at timestamp NOT NULL check ( end_at > start_at ),
    description varchar(100),
    created_by_user bigserial references users(id) on update cascade,
    brackets_limit int default 1 not null check ( brackets_limit > 0 ),
    is_closed boolean default false not null
--     is_private boolean default false NOT NULL
);

select * from tournaments where start_at between now() and  end_at;
-- insert into