import { useCallback, useContext, useEffect, useRef, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { TextareaAutosize } from '@material-ui/core';
import {
  ArrowBackIos,
  Send,
  ImageOutlined,
} from '@material-ui/icons';
import { TYPE_MESSAGE } from "../WsUtils";
import { AuthContext } from "../store/AuthContext";
import { MessageType, WsContext } from "../store/WsContext";
import ImagePicker from "./ImagePicker";
import GalleryModal from "./GalleryModal";
import useOnWsMessage from "../hooks/useOnWsMessage";
import { getChatMessages, getChat, uploadMedia } from "../api/api";

export default function Chat() {
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState<any[]>([]);
  const [text, setText] = useState<string>('');
  const [messages, setMessages] = useState<any[]>([]);
  const [loadingMessages, setLoadingMessages] = useState<boolean>(true);
  const params: any = useParams();
  const history = useHistory();
  const imgInpRef = useRef(null);
  const [open, setOpen] = useState(false);

  const [images, setImages] = useState<any[]>([]);
  const [modalImages, setModalImages] = useState<any[]>([]);

  const { token, userId, logout } = useContext(AuthContext);
  const { sendMessage } = useContext(WsContext);

  const chatId = params.id;

  useOnWsMessage(({ Type, Payload }: MessageType) => {
    if (Type !== 'msg' || Payload.ChatId !== chatId) return;

    setMessages([...messages, Payload]);
  }, [chatId, messages]);

  useEffect(() => {
    getChatMessages(chatId, token)
      .then(result => {
        setMessages([...result])
        setLoadingMessages(false);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        }
      });
  }, [chatId, token]);

  useEffect(() => {
    getChat(params.id, token)
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

  const send = useCallback(() => {
    if (!text?.trim() && images.length <= 0) return;

    const payload = {
      chatId,
      message: text,
      mediaIds: images.map(({ id }) => id)
    };
    sendMessage(TYPE_MESSAGE, payload);

    setText('');
    setImages([]);
  }, [text, images]);

  if (loading) {
    return (
      <div>Loading...</div>
    );
  }

  const onFileChanged = () => {
    const formData = new FormData();
    formData.append('image', (imgInpRef?.current as any)?.files[0]);

    uploadMedia(formData, token)
      .then(response => {
        setImages([...images, {
          id: response.mediaId,
          path: response.mediaPath
        }]);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        }
      });
  };

  return (
    <>
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        height: '100vh'
      }}>
        <div style={{
          padding: 10,
          display: 'flex',
          alignItems: 'center'
        }}>
          <ArrowBackIos style={{ cursor: 'pointer' }} onClick={() => history.push('/chat')}></ArrowBackIos>
          {users.length === 1 && (
            <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }} onClick={() => history.push(`/user/${users[0].Id}`)}>
              <img src={users[0].Image} alt={users[0].Name} style={{ width: 30 }} />
              <span style={{ marginLeft: 5 }}>{users[0].Name}</span>
            </div>
          )}
        </div>
        <div style={{
          flex: 1,
          backgroundColor: '#dedede',
          padding: 5,
          overflow: 'hidden',
          overflowY: 'scroll'
        }}>
          {loadingMessages && <div>Loading...</div>}
          {!loadingMessages && messages.length <= 0 && <div>No history yet.</div>}
          {messages.map(message => (
            <div key={message.Id} style={{
              marginBottom: 5,
              textAlign: message.UserId === userId ? 'right' : 'left'
            }}>
              <div style={{
                padding: 10,
                backgroundColor: message.UserId === userId ? '#9cba9c' : '#cca09f',
                display: 'inline-block',
                borderRadius: 5,
                maxWidth: '70%'
              }}>
                {message.Message}
                {message.Medias?.length > 0 && (
                  <div style={{
                    position: 'relative',
                    marginTop: 10
                  }} onClick={() => {
                    setModalImages(message.Medias);
                    setOpen(true)
                  }}>
                    <img src={message.Medias[0]} style={{
                      width: '100%',
                      borderRadius: 5
                    }} alt={"Message image"} />
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
        <div style={{}}>
          <div style={{
            display: 'flex',
            flexDirection: 'column'
          }}>
            <TextareaAutosize
              minRows={4}
              aria-label="maximum height"
              placeholder="Enter message..."
              style={{ width: '100%', border: '0', display: 'block', outline: 0, padding: 0 }}
              value={text}
              onChange={(e) => setText(e.target.value)}
              onKeyUp={e => {
                if (e.which === 13 && !e.shiftKey) {
                  send();
                }
              }}
            />
            {images.length > 0 && (<ImagePicker
              images={images}
              onAddHandler={() => (imgInpRef?.current as any)?.click()}
              onRemoveHandler={(targetImage: any) => {
                setImages(images.filter(image => image.id !== targetImage.id));
              }}
            ></ImagePicker>)}
            <div style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              width: '100%'
            }}>
              <ImageOutlined
                onClick={() => (imgInpRef?.current as any)?.click()}
                style={{
                  fontSize: 30
                }}
              ></ImageOutlined>
              <input style={{ display: 'none' }} type="file" ref={imgInpRef} onChange={() => onFileChanged()} />
              <Send onClick={() => send()} style={{
                fontSize: 30
              }}></Send>
            </div>
          </div>
        </div>
      </div>
      <GalleryModal open={open} images={modalImages} onClose={() => setOpen(false)}></GalleryModal>
    </>
  );
}
