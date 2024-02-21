CREATE TABLE IF NOT EXISTS invite(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    name VARCHAR(255) NOT NULL,
    greeting VARCHAR(255) NOT NULL,
    lang VARCHAR(2) NOT NULL default 'hu',
    max_adults SMALLINT NOT NULL default 2,
    max_children SMALLINT NOT NULL default 0,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    first_opened_at timestamptz
    last_opened_at timestamptz
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS attendee(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    invite_id uuid NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    age SMALLINT NOT NULL DEFAULT 30,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz
);