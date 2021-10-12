import { Route, Redirect, Switch } from 'react-router-dom';
import Profile from './components/Profile';
import ChatListing from './components/ChatListing';
import Chat from './components/Chat';
import ProfileListing from './components/ProfileListing ';
import Layout from './components/Layout';
import ContactListing from './components/ContactListing';
import Login from './components/Login';
import { useContext } from 'react';
import { AuthContext } from './store/auth-context';
import Signup from './components/Signup';

function App() {
  const { isLoggedIn } = useContext(AuthContext);

  return (
    <div className="App">
      {isLoggedIn && (
        <Switch>
          <Route path="/user/:id" >
            <Layout>
              <Profile />
            </Layout>
          </Route>
          <Route path="/chat" exact>
            <Layout>
              <ChatListing />
            </Layout>
          </Route>
          <Route path="/chat/:id" >
            <Layout>
              <Chat />
            </Layout>
          </Route>
          <Route path="/browse" >
            <Layout>
              <ProfileListing />
            </Layout>
          </Route>
          <Route path="/contacts" >
            <Layout>
              <ContactListing />
            </Layout>
          </Route>
          <Route path="/">
            <Redirect to="/chat" />
          </Route>
        </Switch>
      )}
      {!isLoggedIn && (
        <Switch>
          <Route path="/sign-up" >
            <Signup />
          </Route>
          <Route path="/login" >
            <Login />
          </Route>
          <Route path="/">
            <Redirect to="/login" />
          </Route>
        </Switch>
      )}
    </div>
  );
}

export default App;
