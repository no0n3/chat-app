import { Button, CircularProgress } from "@material-ui/core";
import { useCallback, useContext, useState } from "react";
import { useHistory } from "react-router";
import { addContact, removeContact } from "../api/api";
import { AuthContext } from "../store/AuthContext";

export default function UserItem(props: any) {
  const { user, onAddContact, onRemoveContact } = props;
  const [loading, setLoading] = useState(false);
  const history = useHistory();
  const { token, userId, logout } = useContext(AuthContext);

  const onAddContactC = useCallback(() => {
    if (loading) return;

    setLoading(true);

    addContact(user.Id, token)
      .then(result => {
        setLoading(false);
        onAddContact(user.Id);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setLoading(false);
        }
      });
  }, [user, loading, onAddContact, token]);

  const removeContactC = useCallback(() => {
    if (loading) return;

    setLoading(true);

    removeContact(user.Id, token)
      .then(result => {
        setLoading(false);
        onRemoveContact(user.Id);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setLoading(false);
        }
      });
  }, [user, loading, onRemoveContact, token]);

  const isLoggedUser = user.Id === userId;

  return (
    <div style={{
      padding: 5,
      display: 'flex',
      flexDirection: 'row'
    }}>
      <img src={user.Image} alt={user.Name} style={{ width: 50, cursor: 'pointer' }} onClick={() => history.push(`/user/${user.Id}`)} />
      <div style={{ marginLeft: 5, flex: 1 }}>
        <div style={{ display: 'flex', flexDirection: 'row', cursor: 'pointer' }} onClick={() => history.push(`/user/${user.Id}`)}>
          <h4 style={{ margin: 0 }}>{user.Name}</h4>
          <span style={{ marginLeft: 5 }}>@{user.Username}</span>
        </div>
        <div style={{ marginTop: 5 }}>{user.Description}</div>
      </div>
      {!isLoggedUser && !user.IsContact && (
        <Button onClick={onAddContactC} disabled={loading}>
          {loading && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
          Add</Button>
      )}
      {!isLoggedUser && user.IsContact && (
        <Button onClick={removeContactC} disabled={loading}>
          {loading && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
          Remove
        </Button>
      )}
    </div>
  )
};
