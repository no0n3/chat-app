import { Link, useHistory } from "react-router-dom";
import { Button, CircularProgress } from '@material-ui/core';
import { useContext, useState } from "react";
import { AuthContext } from "../store/auth-context";
import { post } from "../utils/http";

export default function Layout(props: any) {
  const [logingOut, setLogingOut] = useState(false);
  const { token, logout } = useContext(AuthContext);
  const history = useHistory();

  const onLogout = () => {
    if (logingOut) {
      return;
    }
    setLogingOut(true);

    post({
      path: '',
      token,
      payload: {}
    })
      .then(() => {
        logout();
        history.push('/');
      })
      .catch(e => {
        setLogingOut(false);
      });
  };

  return (
    <div style={{
      width: '800px',
      margin: 'auto',
      display: 'flex',
      flexDirection: 'row'
    }}>
      <div style={{ width: 120 }}>
        <div>
          <Link to="/chat" style={{
            textDecoration: 'none',
            color: '#000',
            fontWeight: 500,
            fontSize: 18,
            padding: 5,
            display: 'inline-block'
          }}>Messages</Link>
        </div>
        <div>
          <Link to="/browse" style={{
            textDecoration: 'none',
            color: '#000',
            fontWeight: 500,
            fontSize: 18,
            padding: 5,
            display: 'inline-block'
          }}>Browse</Link>
        </div>
        <div>
          <Link to="/contacts" style={{
            textDecoration: 'none',
            color: '#000',
            fontWeight: 500,
            fontSize: 18,
            padding: 5,
            display: 'inline-block'
          }}>Contacts</Link>
        </div>
        <div>
          <Button onClick={() => onLogout()} disabled={logingOut}>
            {logingOut && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
            Logout
          </Button>
        </div>
      </div>
      <div style={{
        flex: 1,
        borderLeft: '1px solid #000',
        borderRight: '1px solid #000',
        height: '100vh'
      }}>
        {props.children}
      </div>
    </div>
  );
}
