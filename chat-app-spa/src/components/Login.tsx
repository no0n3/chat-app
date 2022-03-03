import { useCallback, useContext, useEffect, useState } from "react";
import { useHistory, Link } from "react-router-dom";
import { TextField, Button, Typography, CircularProgress } from '@material-ui/core';
import { AuthContext } from "../store/AuthContext";
import { isValidEmail } from "../utils/utils";
import { login } from "../api/api";

const FIELDS = ['email', 'password'];

export default function Login() {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [errors, setErrors] = useState<any>({});
  const [touched, setTouched] = useState<any>({});
  const [showErrors, setShowErrors] = useState<any>({});
  const history = useHistory();

  const { loginUser } = useContext(AuthContext);

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
    if (password.length <= 0) {
      errors.password = 'Password must not be empty.';
    }

    setErrors(errors);
  }, [email, password]);

  const onLogin = () => {
    if (loading) return;
    if (Object.keys(errors).some(field => errors[field])) {
      markAsTouched();

      return;
    }

    setLoading(true);

    const payload = {
      email,
      password
    };

    login(payload)
      .then(result => result.data)
      .then((result: any) => {
        const token = result.token;
        const userId = result.userId;
        setEmail('');
        setPassword('');
        setLoading(false);

        loginUser(token, userId);

        history.push('/');
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        setPassword('');
        if (errorData.status === 401) {
          setError('Invalid username or password.')
        } else {
          setError('Internal server error.')
        }
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
        <h1 style={{ margin: 0, textAlign: 'center' }}>Login</h1>
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
                onLogin();
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
                onLogin();
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
        {error && (<div style={{ marginTop: 10, color: 'red', fontSize: 14 }}>{error}</div>)}
        <div>
          <Button
            onClick={() => onLogin()}
            style={{
              marginTop: 15,
              width: '100%'
            }}
            variant="outlined"
            disabled={loading}
          >
            {loading && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
            Login
          </Button>
        </div>
        <div style={{ marginTop: 10 }}>
          <Typography variant="subtitle2" component="div">
            If you don't have a registration you can <Link to="/sign-up">sign up here</Link>
          </Typography>
        </div>
      </div>
    </div>
  );
}
