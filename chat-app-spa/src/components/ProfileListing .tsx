import { Input } from "@material-ui/core";
import { useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext } from "../store/auth-context";
import { get } from "../utils/http";
import UserItem from "./UserItem";

export default function ProfileListing() {
  const [users, setUsers] = useState<any>([]);
  const [filteredUsers, setFilteredUsers] = useState<any>([]);
  const [userFilter, setUserFilter] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const { token, logout } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    if (typeof userFilter !== 'string' || userFilter === '') {
      setFilteredUsers(users);
      return;
    }

    setFilteredUsers(users.filter(({ Name, Username }: { Name: string, Username: string }) => (
      Name.includes(userFilter) || Username.includes(userFilter)
    )));
  }, [userFilter, users]);

  useEffect(() => {
    get('find', token)
      .then(result => {
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

  return (
    <div style={{
      height: '100vh',
      display: 'flex',
      flexDirection: 'column'
    }}>
      <h2 style={{ margin: '5px 10px' }}>Find people</h2>
      <Input
        type="text"
        placeholder="Filter users..."
        style={{ width: '100%' }}
        onChange={(e) => setUserFilter(e.target.value?.trim())}
      ></Input>
      <div style={{
        overflow: 'hidden',
        overflowY: 'scroll',
      }}>
        {loading && (<div>Loading...</div>)}
        {!loading && filteredUsers.length <= 0 && (<div style={{ padding: 10 }}>No users found.</div>)}
        {!loading && filteredUsers.map((user: any) => (<UserItem key={user.Id} user={user}
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
    </div>
  );
}
