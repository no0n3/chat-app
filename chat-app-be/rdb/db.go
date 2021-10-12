package rdb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

var POOL *Pool = &Pool{pool: nil} //initPool()

func initPool() (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(
		context.Background(),
		// "postgres://admin:1234@localhost:5432/chat_app?sslmode=disable",
		"postgresql://postgres:example@db:5432/chat_app",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)

		return nil, err
	}

	return pool, nil
}

func (p *Pool) HandleFunc(f func(*pgxpool.Conn) error) error {
	err := p.handleWithDbConn(f)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (p *Pool) handleWithDbConn(handler func(*pgxpool.Conn) error) error {
	if p.pool == nil {
		pool, err := initPool()
		if err != nil {
			return err
		}

		p.pool = pool
	}
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer conn.Release()

	err = afterConnInit(conn)
	if err != nil {
		return err
	}

	return handler(conn)
}

func afterConnInit(conn *pgxpool.Conn) error {
	_, err := conn.Conn().Prepare(
		context.Background(),
		"select_chats",
		"SELECT chat_id FROM chat_members WHERE user_id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_logged_user_id",
		"SELECT user_id FROM session_tokens WHERE token = $1",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"create_token",
		"INSERT INTO session_tokens (user_id, token, created_at) VALUES ($1::uuid, $2, $3::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_login_info",
		"SELECT id, password FROM users WHERE email = $1",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"create_user",
		"INSERT INTO users (id, name, username, email, password, created_at) VALUES($1::uuid, $2, $3, $4, $5, $6::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_members_chat_id",
		"SELECT chat_id, user_id FROM chat_members WHERE chat_id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_users",
		"SELECT id, name, username, profile_image_id, created_at FROM users WHERE id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_membersx",
		"SELECT id, last_message_id FROM chats WHERE id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_membersxx",
		"SELECT id, chat_id, user_id, message FROM chat_messages WHERE id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_message_medias",
		"SELECT message_id, media_id FROM message_medias WHERE message_id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_users_by_id",
		"SELECT id, name, username, profile_image_id, created_at FROM users WHERE id = any($1)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_members_by_chat",
		"SELECT user_id FROM chat_members WHERE chat_id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_common_chat_id",
		"SELECT cm1.chat_id FROM chat_members AS cm1 JOIN chat_members AS cm2 ON cm1.chat_id = cm2.chat_id WHERE cm1.user_id = $1::uuid AND cm2.user_id = $2::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_chat_messages_by_chat",
		"SELECT id, user_id, message, medias_count, created_at FROM chat_messages WHERE chat_id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"add_chat_message",
		"INSERT INTO chat_messages (id, user_id, chat_id, message, medias_count, created_at) VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5::int, $6::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"add_chat_message_media_rel",
		"INSERT INTO message_medias (message_id, media_id) VALUES ($1::uuid, $2::uuid)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"update_last_message_id_for_chat",
		"UPDATE chats SET last_message_id = $1::uuid WHERE id = $2::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"create_media_metadata",
		"INSERT INTO media_metadatas (id, mime_type, created_at) VALUES ($1::uuid, $2, $3::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"get_media_metadata",
		"SELECT * FROM media_metadatas WHERE id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_user_contacts",
		"SELECT u.id, u.name, u.username, u.description, u.profile_image_id, u.created_at FROM users u JOIN contacts c ON c.added_user_id = u.id WHERE c.adder_user_id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_profile",
		"SELECT id, name, username, profile_image_id, description, created_at FROM users WHERE id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"is_contact",
		"SELECT * FROM contacts WHERE adder_user_id = $1::uuid AND added_user_id = $2::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"select_contacts",
		// "SELECT adder_user_id, added_user_id FROM contacts WHERE (adder_user_id = $1::uuid AND added_user_id = any($2)) OR (added_user_id = $1::uuid AND adder_user_id = any($2))",
		// "SELECT adder_user_id, added_user_id FROM contacts WHERE adder_user_id = $1::uuid OR added_user_id = $1::uuid",
		"SELECT adder_user_id, added_user_id FROM contacts WHERE adder_user_id = $1::uuid OR added_user_id = $1::uuid",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"add_contact",
		"INSERT INTO contacts (adder_user_id, added_user_id, created_at) VALUES ($1::uuid, $2::uuid, $3::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"remove_contact",
		"DELETE FROM contacts WHERE (adder_user_id = $1::uuid AND added_user_id = $2::uuid) OR (added_user_id = $1::uuid AND adder_user_id = $2::uuid)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"create_chat",
		"INSERT INTO chats (id, created_at) VALUES ($1::uuid, $2::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"create_chat_member",
		"INSERT INTO chat_members (chat_id, user_id, created_at) VALUES ($1::uuid, $2::uuid, $3::bigint)",
	)
	if err != nil {
		return err
	}

	_, err = conn.Conn().Prepare(
		context.Background(),
		"change_profile_image",
		"UPDATE users SET profile_image_id = $1::uuid WHERE id = $2::uuid",
	)
	if err != nil {
		return err
	}

	return nil
}
