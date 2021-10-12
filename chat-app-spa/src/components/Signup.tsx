import axios from "axios";
import { useCallback, useContext, useEffect, useState } from "react";
import { useHistory, Link } from "react-router-dom";
import { TextField, Button, Typography, CircularProgress } from '@material-ui/core';
import { AuthContext } from "../store/auth-context";
import { isValidEmail } from "../utils/utils";

const FIELDS = ['email', 'name', 'username', 'password', 'confirmPassword'];

export default function Signup() {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [name, setName] = useState<string>('');
  const [username, setUsername] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [errors, setErrors] = useState<any>({});
  const [touched, setTouched] = useState<any>({});
  const [showErrors, setShowErrors] = useState<any>({});
  const history = useHistory();

  const { login } = useContext(AuthContext);

  const markAsTouched = useCallback(() => {
    const result: any = {};
    FIELDS.forEach(field => result[field] = true);

    setTouched(result);
  }, []);

  useEffect(() => {
    const result: any = {};
    FIELDS.forEach(field => {
      result[field] = !!errors[field] && touched[field];
    });

    setShowErrors(result);
  }, [touched, errors]);

  useEffect(() => {
    const errors: any = {};
    if (!isValidEmail(email)) {
      errors.email = 'Invalid email.';
    }
    if (password.length < 6) {
      errors.password = 'Password must be at least 6 characters';
    }
    if (confirmPassword !== password) {
      errors.confirmPassword = 'Passwords do not match.';
    }
    if (name.length < 3) {
      errors.name = 'Name must be at least 3 characters.';
    }
    if (username.length < 3) {
      errors.username = 'Username must be at least 3 characters.';
    }

    setErrors(errors);
  }, [email, name, username, password, confirmPassword]);

  const onSignup = () => {
    if (loading) {
      return;
    }
    if (Object.keys(errors).some(field => errors[field])) {
      markAsTouched();

      return;
    }

    setLoading(true);

    const payload = {
      email,
      password,
      name,
      username
    };

    axios.post(`${process.env.REACT_APP_ENDPOINT}/api/sign-up`, payload)
      .then(result => result.data)
      .then((result: any) => {
        const token = result.token;
        const userId = result.userId;

        setEmail('');
        setPassword('');
        setConfirmPassword('');
        setLoading(false);

        login(token, userId);

        history.push('/');
      })
      .catch((e) => {
        setPassword('');
        setConfirmPassword('');
        setError('Internal server error.')
        setLoading(false);
      });
  };

  return (
    <div>
      <div style={{
        width: 400,
        margin: 'auto',
        marginTop: 100,
        padding: 30,
        border: '1px solid #000',
        borderRadius: 25
      }}>
        <h1 style={{ margin: 0, textAlign: 'center' }}>Sign up</h1>
        <div>
          <TextField
            error={showErrors.email}
            type="email"
            label="Email"
            onFocus={() => setTouched({ ...touched, email: true })}
            onChange={(e) => {
              setEmail(e.target.value);
              setError('');
            }}
            onKeyUp={(e) => {
              if (e.which === 13) {
                onSignup();
              }
            }}
            value={email}
            style={{
              width: '100%'
            }}
          />
        </div>
        {showErrors.email && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{errors.email}</div>)}

        <div>
          <TextField
            error={showErrors.name}
            type="text"
            label="Name"
            onFocus={() => setTouched({ ...touched, name: true })}
            onChange={(e) => {
              setName(e.target.value);
              setError('');
            }}
            onKeyUp={(e) => {
              if (e.which === 13) {
                onSignup();
              }
            }}
            value={name}
            style={{
              marginTop: 15,
              width: '100%'
            }}
          />
        </div>
        {showErrors.name && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{errors.name}</div>)}

        <div>
          <TextField
            error={showErrors.username}
            type="text"
            label="Username"
            onFocus={() => setTouched({ ...touched, username: true })}
            onChange={(e) => {
              setUsername(e.target.value);
              setError('');
            }}
            onKeyUp={(e) => {
              if (e.which === 13) {
                onSignup();
              }
            }}
            value={username}
            style={{
              marginTop: 15,
              width: '100%'
            }}
          />
        </div>
        {showErrors.username && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{errors.username}</div>)}

        <div>
          <TextField
            error={showErrors.password}
            label="Password"
            type="password"
            onFocus={() => setTouched({ ...touched, password: true })}
            onChange={(e) => {
              setPassword(e.target.value);
              setError('');
            }}
            onKeyUp={(e) => {
              if (e.which === 13) {
                onSignup();
              }
            }}
            value={password}
            style={{
              marginTop: 15,
              width: '100%'
            }}
          />
        </div>
        {showErrors.password && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{errors.password}</div>)}

        <div>
          <TextField
            error={showErrors.confirmPassword}
            label="Confirm password"
            type="password"
            onFocus={() => setTouched({ ...touched, confirmPassword: true })}
            onChange={(e) => {
              setConfirmPassword(e.target.value);
              setError('');
            }}
            onKeyUp={(e) => {
              if (e.which === 13) {
                onSignup();
              }
            }}
            value={confirmPassword}
            style={{
              marginTop: 15,
              width: '100%'
            }}
          />
        </div>
        {showErrors.confirmPassword && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{errors.confirmPassword}</div>)}

        {error && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{error}</div>)}
        <div>
          <Button
            onClick={() => onSignup()}
            style={{
              marginTop: 15,
              width: '100%'
            }}
            variant="outlined"
            disabled={loading}
          >
            {loading && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
            Sign up
          </Button>
        </div>
        <div style={{ marginTop: 15 }}>
          <Typography variant="subtitle2" component="div">
            If you already have a registration you can <Link to="/login">login here</Link>
          </Typography>
        </div>
      </div>
    </div>
  );
}
