import { useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { getContacts } from "../api/api";
import { AuthContext } from "../store/AuthContext";
import UserItem from "./UserItem";

function useGetUsers(): [any[], (users: any[]) => void, boolean] {
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState<any[]>([]);

  const { token, logout } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    getContacts(token)
      .then(result => {
        setUsers(result);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setUsers([]);
        }
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return [users, setUsers, loading];
}

export default function ContactListing() {
  const [users, setUsers, loading] = useGetUsers();

  return (
    <div style={{
      overflow: 'hidden',
      overflowY: 'scroll',
      height: '100vh'
    }}>
      <h2 style={{ margin: '5px 10px' }}>My contacts</h2>
      {loading && (<div>Loading...</div>)}
      {!loading && users.length <= 0 && (<div style={{ padding: 10 }}>No contacts found.</div>)}
      {!loading && users.map(user => (<UserItem key={user.Id} user={user}
        onAddContact={(userId: string) => {
          setUsers(users.map((user: any) => {
            if (user.Id === userId) {
              return {
                ...user,
                IsContact: true
              };
            }

            return user;
          }));
        }}
        onRemoveContact={(userId: string) => {
          setUsers(users.map((user: any) => {
            if (user.Id === userId) {
              return {
                ...user,
                IsContact: false
              };
            }

            return user;
          }));
        }}
      ></UserItem>))}
    </div>
  );
}
