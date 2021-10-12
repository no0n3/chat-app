import { Button, CircularProgress } from "@material-ui/core";
import Cropper from "cropperjs";
import { useContext, useEffect, useRef, useState } from "react";
import { useHistory } from "react-router";
import { AuthContext } from "../store/auth-context";
import { post } from "../utils/http";

export default function CustCropper(props: any) {
  const [cropper, setCropper] = useState<Cropper | null>(null);
  const [saving, setSaving] = useState(false);
  const { token, logout } = useContext(AuthContext);
  const cropperImageRef = useRef<any>(null);
  const history = useHistory();
  const { src, closeModal } = props;

  const handleSave = () => {
    if (saving) {
      return;
    }
    setSaving(true);

    const canvas = cropper?.getCroppedCanvas();

    const file = dataURLtoBlob(canvas?.toDataURL());
    const formData = new FormData;
    formData.append('image', file);

    post({
      path: 'upload-profile-image',
      token,
      payload: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
      .then(response => {
        closeModal(response.mediaPath);
      })
      .catch((e) => {
        const errorData = e?.toJSON();
        if (errorData.status === 401) {
          logout();
          history.push('/login');
        } else {
          setSaving(false);
        }
      });
  }

  const dataURLtoBlob = (dataURL: any) => {
    let array, binary, i, len;
    binary = atob(dataURL.split(',')[1]);
    array = [];
    i = 0;
    len = binary.length;
    while (i < len) {
      array.push(binary.charCodeAt(i));
      i++;
    }

    return new Blob([new Uint8Array(array)], {
      type: 'image/png'
    });
  };

  useEffect(() => {
    if (!src) {
      return;
    }

    setCropper(new Cropper(cropperImageRef.current as any, {
      zoomable: false,
      aspectRatio: 1,
    }));
  }, [src]);

  if (!src) {
    return (<></>);
  }

  return (
    <>
      <div style={{ position: 'relative' }}>
        <img
          ref={cropperImageRef}
          src={src}
          style={{ width: '100%', display: 'block' }}
        />
      </div>
      <div style={{
        padding: 5,
        display: 'flex',
        flexDirection: 'row-reverse'
      }}>
        <Button onClick={() => handleSave()} disabled={saving}>
          {saving && <CircularProgress size={20} style={{ marginRight: 5 }}></CircularProgress>}
          Save
        </Button>
      </div>
    </>
  );
}
