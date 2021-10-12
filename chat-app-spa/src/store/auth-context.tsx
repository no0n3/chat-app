import { createContext, useEffect, useState } from "react";

export const AuthContext = createContext({
  token: '',
  userId: '',
  isLoggedIn: false,
  login: (token: string, userId: string) => {},
  logout: () => {},
} as any);

export default function AuthContextProvider(props: any) {
  const [token, setToken] = useState(localStorage.getItem('ca_token') as string);
  const [userId, setUserId] = useState(localStorage.getItem('ca_user_id') as string);

  useEffect(() => {
    if (!token) {
      localStorage.removeItem('ca_token');
    } else {
      localStorage.setItem('ca_token', token);
    }
  }, [token]);

  useEffect(() => {
    if (!userId) {
      localStorage.removeItem('ca_user_id');
    } else {
      localStorage.setItem('ca_user_id', userId);
    }
  }, [userId]);

  return (
    <AuthContext.Provider value={{
      token,
      userId,
      isLoggedIn: !!token,
      login: (newToken: string, newUserId: string) => {
        setToken(newToken);
        setUserId(newUserId);
      },
      logout: () => {
        setToken('');
        setUserId('');
      },
    }}>
      {props.children}
    </AuthContext.Provider>
  );
};
