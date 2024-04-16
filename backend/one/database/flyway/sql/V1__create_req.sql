CREATE TABLE request (
    req_id text NOT NULL,
    user_id text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    req jsonb NOT NULL,
    PRIMARY KEY (req_id)
);