CREATE EXTENSION pgcrypto;
CREATE EXTENSION citext;

CREATE TABLE "staff_role" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE "staff" (
    id SERIAL PRIMARY KEY,
    staff_role_id INTEGER NOT NULL,
    username VARCHAR(15) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    CONSTRAINT fk_staff_role FOREIGN KEY (staff_role_id) REFERENCES "staff_role" (id)
);

CREATE TABLE "staff_action" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) UNIQUE NOT NULL
);

CREATE TABLE "staff_permission" (
    id SERIAL PRIMARY KEY,
    staff_role_id INTEGER NOT NULL,
    staff_action_id INTEGER NOT NULL,
    UNIQUE (staff_role_id, staff_action_id),
    CONSTRAINT fk_staff_allowed_role FOREIGN KEY (staff_role_id) REFERENCES "staff_role" (id),
    CONSTRAINT fk_staff_allowed_possible FOREIGN KEY (staff_action_id) REFERENCES "staff_action" (id)
);

CREATE TABLE "staff_log" (
    id SERIAL PRIMARY KEY,
    staff_id INTEGER NOT NULL,
    staff_action_id INTEGER NOT NULL,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    ip_address INET NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_staff_log_staff FOREIGN KEY (staff_id) REFERENCES "staff" (id),
    CONSTRAINT fk_staff_log_action FOREIGN KEY (staff_action_id) REFERENCES "staff_action" (id)
);

INSERT INTO staff_role (name) VALUES
    ('Developer - General'), ('Developer - Game'), ('Developer - Web'), ('Developer - DB/API'),
    ('Artist - General'), ('Artist - Character'), ('Artist - Animator'), ('Artist - Environment'),
    ('Designer - General'), ('Designer - Environment'), ('Designer - Gameplay'), ('Designer - Level'),
    ('Designer - UI/UX'), ('Designer - Character'), ('Writer'), ('Management'), ('Producer'), ('HR'), ('Sales'),
    ('Marketing'), ('Web - Moderation'), ('Web - Support'), ('QA - General'), ('QA - Game testing'),
    ('Sound - Effects'), ('Sound - Music'), ('Sound - Voice'), ('Translator');

INSERT INTO staff (staff_role_id, username, password, email) 
VALUES (1,'xah', crypt('xConfused9', gen_salt('bf')), 'xahkun@gmail.com');

INSERT into staff_action (name) VALUES ('staff-create'), ('staff-update'), ('staffrole-create'),
('staffrole-update'), ('stafflog-read'), ('staffaction-create'), ('staffaction-update'), ('staff-action-delete'),
('staffpermission-create'), ('staffpermission-update'), ('staffpermission-delete');

INSERT INTO staff_permission (staff_role_id, staff_action_id) VALUES (1,1), (1,2), (1,3), (1,4), (1,5),
(1,6), (1,7), (1,8), (1,9), (1,10), (1,11);

-- Create player table: This table records all users (players) of the game
CREATE TABLE "player" (
  id SERIAL PRIMARY KEY,
  username VARCHAR(15) NOT NULL,
  password TEXT NOT NULL,
  email CITEXT UNIQUE,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  verified_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  disabled_at TIMESTAMP
);

-- Create login history table (for tracking purposes)
CREATE TABLE "login_history" (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL,
    login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logout_time TIMESTAMP,
    session_data TEXT,
    CONSTRAINT fk_loginhistory_player FOREIGN KEY (player_id) REFERENCES "player" (id)
);

-- Create Guilds
CREATE TABLE "guild" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    disbanded_at TIMESTAMP
);

CREATE TABLE "guild_history" (
    id SERIAL PRIMARY KEY,
    guild_id INTEGER NOT NULL,
    icon TEXT NOT NULL,
    tagline VARCHAR(50),
    comments TEXT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_guild_history_guild FOREIGN KEY (guild_id) REFERENCES "guild" (id)
);

CREATE TABLE "guild_role" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE "guild_action_possible" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE "guild_action_allowed" (
    id SERIAL PRIMARY KEY,
    guild_role_id INTEGER NOT NULL,
    guild_action_id INTEGER NOT NULL,
    CONSTRAINT fk_guild_allowed_role FOREIGN KEY (guild_role_id) REFERENCES "guild_role" (id),
    CONSTRAINT fk_guild_allowed_possible FOREIGN KEY (guild_action_id) REFERENCES "guild_action_possible" (id)
);

CREATE TABLE "guild_member" (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL,
    guild_id INTEGER NOT NULL,
    guild_role_id INTEGER NOT NULL,
    tagline VARCHAR(25),
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    CONSTRAINT fk_guild_member_player FOREIGN KEY (player_id) REFERENCES "player" (id),
    CONSTRAINT fk_guild_member_guild FOREIGN KEY (guild_id) REFERENCES "guild" (id),
    CONSTRAINT fk_guild_member_role FOREIGN KEY (guild_role_id) REFERENCES "guild_role" (id)
);

CREATE TABLE "guild_member_history" (
    id SERIAL PRIMARY KEY,
    guild_member_id INTEGER,
    guild_id INTEGER NOT NULL,
    guild_role_id INTEGER NOT NULL,
    tagline VARCHAR(25),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_guild_memberhistory_member FOREIGN KEY (guild_member_id) REFERENCES "guild_member" (id),
    CONSTRAINT fk_guild_memberhistory_guild FOREIGN KEY (guild_id) REFERENCES "guild" (id),
    CONSTRAINT fk_guild_memberhistory_role FOREIGN KEY (guild_role_id) REFERENCES "guild_role" (id)
);

-- Create hero tables
CREATE TABLE "hero" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15),
    class INTEGER NOT NULL,
    rarity INTEGER NOT NULL,
    element INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "hero_background" (
    hero_id INTEGER PRIMARY KEY,
    short_description TEXT NOT NULL,
    long_description TEXT NOT NULL,
    CONSTRAINT fk_hero_background_hero FOREIGN KEY (hero_id) REFERENCES "hero" (id)
);

CREATE TABLE "hero_base_stat" (
    hero_id INTEGER PRIMARY KEY,
    attack INTEGER NOT NULL,
    crit_chance INTEGER NOT NULL,
    crit_damage INTEGER NOT NULL,
    hit_chance INTEGER NOT NULL,
    effectiveness INTEGER NOT NULL,
    health INTEGER NOT NULL,
    defence INTEGER NOT NULL,
    evasion INTEGER NOT NULL,
    resistance INTEGER NOT NULL,
    CONSTRAINT fk_hero_basestat_hero FOREIGN KEY (hero_id) REFERENCES "hero" (id)
);

CREATE TABLE "hero_skill" (
    id SERIAL PRIMARY KEY,
    hero_id INTEGER NOT NULL,
    icon TEXT UNIQUE NOT NULL,
    description TEXT UNIQUE NOT NULL,
    is_active_skill BOOLEAN DEFAULT TRUE,
    CONSTRAINT fk_hero_skill_hero FOREIGN KEY (hero_id) REFERENCES "hero" (id)
);

CREATE UNIQUE INDEX "hero_skill_id" ON hero_skill (id, hero_id);

CREATE TABLE "hero_skill_multiplier" (
    hero_skill_id INTEGER NOT NULL,
    hero_id INTEGER NOT NULL,
    power DECIMAL NOT NULL,
    attack DECIMAL NOT NULL,
    crit_damage DECIMAL NOT NULL DEFAULT 0,
    defence DECIMAL NOT NULL DEFAULT 0,
    health DECIMAL NOT NULL DEFAULT 0,
    percent_health DECIMAL NOT NULL DEFAULT 0,
    vs_debuffed DECIMAL NOT NULL DEFAULT 0,
    vs_evaded DECIMAL NOT NULL DEFAULT 0,
    vs_high_hp DECIMAL NOT NULL DEFAULT 0,
    vs_high_atk DECIMAL NOT NULL DEFAULT 0,
    CONSTRAINT fk_hero_skill_multiplier FOREIGN KEY (hero_skill_id, hero_id) REFERENCES "hero_skill" (id, hero_id)
);

CREATE TABLE "hero_action_possible" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE "hero_action_allowed" (
    id SERIAL PRIMARY KEY,
    execution_order INTEGER NOT NULL,
    hero_skill_id INTEGER NOT NULL,
    hero_action_id INTEGER NOT NULL,
    chance INTEGER NOT NULL,
    CONSTRAINT fk_hero_action_allowed_skill FOREIGN KEY (hero_skill_id) REFERENCES "hero_skill" (id),
    CONSTRAINT fk_hero_action_allowed_possible FOREIGN KEY (hero_action_id) REFERENCES "hero_action_possible" (id)
);

-- Create tables for other resources
CREATE TABLE "currency" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ring_skill" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE "ring" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    skill_id INTEGER NOT NULL,
    CONSTRAINT fk_ring_skill FOREIGN KEY (skill_id) REFERENCES "ring_skill" (id)
);

CREATE TABLE "hero_banner" (
    id SERIAL PRIMARY KEY,
    hero_id INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    run_start TIMESTAMP,
    run_end TIMESTAMP,
    CONSTRAINT fk_hero_banner_hero FOREIGN KEY (hero_id) REFERENCES "hero" (id)
);

-- Add player heroes & resources
CREATE TABLE "player_currency" (
    player_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    CONSTRAINT fk_currency_player FOREIGN KEY (player_id) REFERENCES "player" (id)
);

CREATE UNIQUE INDEX "player_currency_id" ON player_currency (player_id, currency_id);

CREATE TABLE "player_transaction" (
    id SERiAL PRIMARY KEY,
    player_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    change INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_player_currency_transaction FOREIGN KEY (player_id, currency_id) REFERENCES "player_currency" (player_id, currency_id)
);

CREATE UNIQUE INDEX "player_transaction_id" ON player_transaction (id, player_id);

CREATE TABLE "summon" (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    hero_id INTEGER NOT NULL,
    hero_banner_id INTEGER NOT NULL,
    level INTEGER DEFAULT 1,
    friendship DECIMAL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_summon_transaction FOREIGN KEY (transaction_id, player_id) REFERENCES "player_transaction" (id, player_id),
    CONSTRAINT fk_summon_hero FOREIGN KEY (hero_id) REFERENCES "hero" (id),
    CONSTRAINT fk_summon_banner FOREIGN KEY (hero_banner_id) REFERENCES "hero_banner" (id)
);

CREATE TABLE "summon_stats" {
    summon_id INTEGER NOT NULL UNIQUE,
    attack INTEGER NOT NULL,
    crit_chance INTEGER NOT NULL,
    crit_damage INTEGER NOT NULL,
    hit_chance INTEGER NOT NULL,
    effectiveness INTEGER NOT NULL,
    health INTEGER NOT NULL,
    defence INTEGER NOT NULL,
    evasion INTEGER NOT NULL,
    resistance INTEGER NOT NULL,
    constraint fk_summon_stats FOREIGN KEY (summon_id) REFERENCES "summon" (id)
}

CREATE TABLE "player_ring" (
    player_id INTEGER NOT NULL,
    ring_id INTEGER NOT NULL UNIQUE,
    summon_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_player_ring FOREIGN KEY (ring_id) REFERENCES "ring" (id),
    CONSTRAINT fk_player_ring_hero FOREIGN KEY (summon_id) REFERENCES "summon" (id)
);