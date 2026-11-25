CREATE SCHEMA graphics;

--учетные записи и роли
CREATE TABLE graphics.users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 3 AND 100),
    email VARCHAR(254) NOT NULL UNIQUE CHECK (char_length(email) > 0),
    password_hash TEXT NOT NULL CHECK (char_length(password_hash) > 0),
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'author', 'admin')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version BIGINT NOT NULL DEFAULT 1 CHECK (version > 0)
);
--разделы каталога
CREATE TABLE graphics.categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(name) BETWEEN 1 AND 100),
    version BIGINT NOT NULL DEFAULT 1 CHECK (version > 0)
);
--метки для поиска и фильтрации
CREATE TABLE graphics.tags (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(name) BETWEEN 1 AND 100)
);
--карточки
CREATE TABLE graphics.resources (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES graphics.users (id),
    category_id UUID NOT NULL REFERENCES graphics.categories (id),
    title VARCHAR(100) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'draft'
                             CHECK (status IN (
                                 'draft', 'pending', 'published',
                                 'rejected', 'hidden', 'archived'
                             )),
    version BIGINT NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (updated_at >= created_at)
);
--ресурсы для карточки
CREATE TABLE graphics.resource_files (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    kind TEXT NOT NULL CHECK (kind IN ('original', 'image', 'preview_3d')),
    storage_key TEXT NOT NULL UNIQUE CHECK (char_length(storage_key) > 0),
    original_name TEXT NOT NULL CHECK (char_length(original_name) > 0),
    format TEXT NOT NULL CHECK (char_length(format) > 0),
    size_bytes BIGINT NOT NULL CHECK (size_bytes > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    retired_at TIMESTAMPTZ,

    CHECK (retired_at IS NULL OR retired_at >= created_at)
);
--теги карточки
CREATE TABLE graphics.resource_tags (
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    tag_id UUID NOT NULL REFERENCES graphics.tags (id),

    PRIMARY KEY (resource_id, tag_id)
);
--избранное
CREATE TABLE graphics.favorites (
    user_id UUID NOT NULL REFERENCES graphics.users (id),
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, resource_id)
);
--решения проверки
CREATE TABLE graphics.moderation_reviews (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    moderator_id UUID NOT NULL REFERENCES graphics.users (id),
    reviewed_version BIGINT NOT NULL CHECK (reviewed_version > 0),
    decision TEXT NOT NULL CHECK (decision IN ('approved', 'rejected')),
    comment TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (resource_id, reviewed_version),
    CHECK (decision <> 'rejected' OR char_length(comment) > 0)
);
--история изменений карточек и файлов
CREATE TABLE graphics.resource_events (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    actor_id UUID NOT NULL REFERENCES graphics.users (id),
    action TEXT NOT NULL CHECK (char_length(action) > 0),
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--просмотры опубликованных карточек
CREATE TABLE graphics.resource_views (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES graphics.resources (id),
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--история выдачи оригиналов и источник статистики
CREATE TABLE graphics.downloads (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES graphics.users (id),
    file_id UUID NOT NULL REFERENCES graphics.resource_files (id),
    issued_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
