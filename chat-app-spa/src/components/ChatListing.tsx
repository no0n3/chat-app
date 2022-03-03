import { useContext, useEffect, useRef, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext } from "../store/AuthContext";
import { Input } from '@material-ui/core';
import { Reply } from '@material-ui/icons';
import useOnWsMessage from "../hooks/useOnWsMessage";
import { getChats } from "../api/api";
import { MessageType } from "../store/WsContext";

export default function ChatListing() {
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState<any[]>([]);
  const { token, logout } = useContext(AuthContext);
  const history = useHistory();

  const usersRef = useRef<any[]>();

  useOnWsMessage(({ Type, Payload }: MessageType) => {
    if (Type !== 'msg') return;

    if (!Payload || !usersRef.current) return;

    const newUsers: any[] = [];
    usersRef.current.forEach(user => {
      if (user.ChatId === Payload.ChatId) {
        newUsers.push({
          ...user,
          LastMessage: { ...Payload }
        });
      } else {
        newUsers.push({ ...user });
      }
    });

    usersRef.current = newUsers;
    setUsers(newUsers);
  });

  const [filteredUsers, setFilteredUsers] = useState<any>([]);
  const [userFilter, setUserFilter] = useState<string>('');

  useEffect(() => {
    if (typeof userFilter !== 'string' || userFilter === '') {
      setFilteredUsers(users);
      return;
    }

    const tu = users.filter(({ Name, Username }: { Name: string, Username: string }) => (
      Name.includes(userFilter) || Username.includes(userFilter)
    ));

    setFilteredUsers(tu);
  }, [userFilter, users]);

  useEffect(() => {
    getChats(token)
      .then(result => {
        usersRef.current = result;
        setUsers(result);
        setLoading(false);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        }
      });
  }, []);

  if (loading) {
    return (
      <div>Loading...</div>
    );
  }

  return (
    <div style={{
      height: '100vh',
      display: 'flex',
      flexDirection: 'column'
    }}>
      <Input
        type="text"
        placeholder="Filter users..."
        style={{ width: '100%', padding: 5 }}
        onChange={(e: any) => setUserFilter(e.target.value?.trim())}
      ></Input>
      <div style={{
        flex: 1,
        overflow: 'hidden',
        overflowY: 'scroll',
      }}>
        {!loading && filteredUsers.length <= 0 && <div style={{ padding: 10 }}>No contacts found.</div>}
        {filteredUsers.map((user: any) => <ChatMemberItem key={user.Id} user={user}></ChatMemberItem>)}
      </div>
    </div>
  );
}

function ChatMemberItem(props: any) {
  const { userId } = useContext(AuthContext);
  const history = useHistory();
  const { user } = props;

  return (
    <div key={user.Id} onClick={() => history.push(`/chat/${user.ChatId}`)} style={{
      display: 'flex',
      padding: 5,
      cursor: 'pointer'
    }}>
      <img src={user.Image} alt={user.Name} style={{ width: 50 }} />
      <div style={{
        marginLeft: 5
      }}>
        <span style={{ fontSize: 18 }}>{user.Name}</span>
        <span style={{ fontWeight: 700, marginLeft: 5 }}>@{user.Username}</span>
        {user.LastMessage && (<div style={{
          display: 'flex',
          alignItems: 'center'
        }}>
          {userId === user.LastMessage.UserId && <Reply></Reply>}
          <div style={{
            marginLeft: 5,
            fontSize: 15,
            color: '#828282s'
          }}>{user.LastMessage.Message}</div>
        </div>)}
      </div>
    </div>
  );
}
