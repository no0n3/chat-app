import { useCallback, useContext, useEffect, useRef, useState } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { AuthContext } from '../store/AuthContext';
import { Button, CircularProgress, Dialog } from '@material-ui/core';
import CustCropper from './Cropper';
import { addContact, getUser, removeContact } from '../api/api';

export default function Profile() {
  const [user, setUser] = useState<any>({});
  const [loading, setLoading] = useState(true);
  const [loadingContact, setLoadingContact] = useState(false);
  const [open, setOpen] = useState(false);
  const { token, userId, logout } = useContext(AuthContext);
  const ref = useRef(null);
  const history = useHistory();

  const [src, setSrc] = useState<string>('');

  const params: any = useParams();

  useEffect(() => {
    getUser(params.id, token)
      .then(result => {
        setUser(result);
        setLoading(false);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else if (errorData.status === 404) {
          history.push('/');
        }
      });
  }, []);

  const onAddContact = useCallback(() => {
    if (loadingContact) return;

    setLoadingContact(true);

    addContact(user.Id, token)
      .then(result => {
        setLoadingContact(false);
        setUser({
          ...user,
          IsContact: true
        });
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setLoadingContact(false);
        }
      });
  }, [user, loadingContact]);

  const onRemoveContact = useCallback(() => {
    if (loadingContact) return;

    setLoadingContact(true);

    removeContact(user.Id, token)
      .then(result => {
        setLoadingContact(false);
        setUser({
          ...user,
          IsContact: false
        });
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setLoadingContact(false);
        }
      });
  }, [user, loadingContact]);

  if (loading) {
    return (<div>Loading...</div>);
  }

  const getImage = (image: any) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        resolve(e);
      };

      reader.readAsDataURL(image);
    });
  };

  const openUpload = () => {
    if (userId !== user.Id) return;
    const inp = ref?.current as any;
    if (!inp) return;

    inp.click();
  };

  const onFileChange = () => {
    const inp = ref?.current as any;
    if (!inp) return;

    const image = inp.files[0];

    setOpen(true);
    getImage(image)
      .then((e: any) => {
        setSrc(e.target.result as string);
      });
  };

  const isLoggedUser: boolean = userId === user.Id;

  return (
    <>
      <div style={{
        display: 'flex',
        flexDirection: 'row',
        padding: 5
      }}>
        <div style={{
          borderRadius: '100%',
          overflow: 'hidden',
          width: 100,
          height: 100,
          border: '1px solid #000',
          position: 'relative'
        }}>
          <img src={user.Image} alt={user.Name} style={{ width: 100, display: 'block' }} crossOrigin="anonymous" />
          {isLoggedUser && (<div onClick={openUpload} style={{
            backgroundColor: 'rgba(0, 0, 0, .5)',
            color: '#fff',
            textAlign: 'center',
            position: 'absolute',
            bottom: 0,
          }}>
            Change image
            <input style={{ display: 'none' }} type="file" accept="image/*" ref={ref} onChange={() => onFileChange()} />
          </div>)}
        </div>
        <div style={{
          marginLeft: 10
        }}>
          <div style={{ fontSize: 20 }}>{user.Name}</div>
          <div>@{user.Username}</div>
          <div style={{ paddingTop: 10 }}>{user.Description ? user.Description : 'No description yet.'}</div>

          <div>
            {!user.IsContact && !isLoggedUser && (
              <Button onClick={onAddContact} disabled={loadingContact}>
                {loadingContact && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
                Add contact
              </Button>
            )}
            {user.IsContact && !isLoggedUser && (
              <Button onClick={onRemoveContact} disabled={loadingContact}>
                {loadingContact && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
                Remove contact
              </Button>
            )}
          </div>
        </div>
      </div>

      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        aria-labelledby="parent-modal-title"
        aria-describedby="parent-modal-description"
      >
        <div style={{
          width: 500
        }}>
          <CustCropper src={src} closeModal={(mediaPath: string) => {
            setOpen(false);
            setUser({
              ...user,
              Image: mediaPath
            })
          }}></CustCropper>
        </div>
      </Dialog>
    </>
  );
}
