CREATE TABLE media_metadatas (
  id UUID PRIMARY KEY,
  mime_type VARCHAR(255) NOT NULL,
  created_at BIGINT NOT NULL
);

CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  description TEXT,
  profile_image_id UUID REFERENCES media_metadatas(id),
  last_online_at BIGINT,
  created_at BIGINT NOT NULL
);

CREATE TABLE contacts (
  adder_user_id UUID NOT NULL REFERENCES users(id),
  added_user_id UUID NOT NULL REFERENCES users(id),
  created_at BIGINT NOT NULL
);

CREATE TABLE session_tokens (
  token VARCHAR(255) PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  created_at BIGINT NOT NULL
);

CREATE TABLE chats (
  id UUID PRIMARY KEY,
  created_at BIGINT NOT NULL
);

CREATE TABLE chat_members (
  chat_id UUID NOT NULL REFERENCES chats(id),
  user_id UUID NOT NULL REFERENCES users(id),
  created_at BIGINT NOT NULL
);

CREATE TABLE chat_messages (
  id UUID PRIMARY KEY,
  chat_id UUID NOT NULL REFERENCES chats(id),
  user_id UUID NOT NULL REFERENCES users(id),
  message TEXT,
  medias_count INTEGER NOT NULL,
  created_at BIGINT NOT NULL
);

ALTER TABLE chats ADD COLUMN last_message_id UUID REFERENCES chat_messages(id);

CREATE TABLE message_medias (
  message_id UUID NOT NULL REFERENCES chat_messages(id),
  media_id UUID NOT NULL REFERENCES media_metadatas(id)
  -- created_at BIGINT NOT NULL
);
