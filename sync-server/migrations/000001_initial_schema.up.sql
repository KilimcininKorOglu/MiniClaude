create extension if not exists pgcrypto;

create table users (
    id uuid primary key default gen_random_uuid(),
    email text not null unique,
    username text not null unique,
    password_hash text not null,
    status text not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table workspaces (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    slug text not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table workspace_members (
    workspace_id uuid not null references workspaces(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    role text not null,
    created_at timestamptz not null default now(),
    primary key (workspace_id, user_id)
);

create table web_sessions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    workspace_id uuid references workspaces(id) on delete set null,
    jti text not null unique,
    expires_at timestamptz not null,
    revoked_at timestamptz,
    created_at timestamptz not null default now()
);

create table settings_documents (
    workspace_id uuid primary key references workspaces(id) on delete cascade,
    document jsonb not null default '{}'::jsonb,
    version bigint not null default 1,
    checksum text not null,
    updated_by uuid references users(id) on delete set null,
    updated_at timestamptz not null default now()
);

create table settings_events (
    id bigserial primary key,
    workspace_id uuid not null references workspaces(id) on delete cascade,
    version bigint not null,
    event_type text not null,
    payload jsonb not null,
    created_by uuid references users(id) on delete set null,
    created_at timestamptz not null default now(),
    unique (workspace_id, version)
);

create table provider_configs (
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces(id) on delete cascade,
    name text not null,
    base_url text,
    model text,
    default_models jsonb not null default '{}'::jsonb,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (workspace_id, name)
);

create table provider_secrets (
    provider_id uuid primary key references provider_configs(id) on delete cascade,
    key_version integer not null default 1,
    nonce bytea not null,
    ciphertext bytea not null,
    updated_at timestamptz not null default now()
);

create table clients (
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    name text not null,
    machine_id text,
    created_at timestamptz not null default now(),
    last_seen_at timestamptz,
    revoked_at timestamptz
);

create table client_sessions (
    id uuid primary key default gen_random_uuid(),
    client_id uuid not null references clients(id) on delete cascade,
    workspace_id uuid not null references workspaces(id) on delete cascade,
    process_id text,
    status text not null default 'active',
    connected_at timestamptz not null default now(),
    last_seen_at timestamptz not null default now(),
    terminated_at timestamptz
);

create table device_login_requests (
    id uuid primary key default gen_random_uuid(),
    device_code text not null unique,
    user_code text not null unique,
    workspace_id uuid references workspaces(id) on delete cascade,
    approved_by uuid references users(id) on delete set null,
    client_id uuid references clients(id) on delete set null,
    status text not null default 'pending',
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    approved_at timestamptz
);

create table audit_events (
    id bigserial primary key,
    workspace_id uuid references workspaces(id) on delete cascade,
    user_id uuid references users(id) on delete set null,
    actor_type text not null,
    event_type text not null,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now()
);

create index idx_settings_events_workspace_version on settings_events (workspace_id, version);
create index idx_client_sessions_workspace_status on client_sessions (workspace_id, status);
create index idx_device_login_requests_user_code on device_login_requests (user_code);
create index idx_audit_events_workspace_created_at on audit_events (workspace_id, created_at desc);
