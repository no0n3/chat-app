import { Dialog } from "@material-ui/core";
import {
  HighlightOffOutlined,
  KeyboardArrowLeft,
  KeyboardArrowRight
} from "@material-ui/icons";
import { useEffect, useState } from "react";

export default function GalleryModal(props: any) {
  const [position, setPosition] = useState(0);
  const { open, images, onClose } = props;

  useEffect(() => {
    setPosition(0);
  }, [images])

  return (
    <Dialog
      open={open}
      aria-labelledby="parent-modal-title"
      aria-describedby="parent-modal-description"
    >
      <div style={{
        width: 500,
        position: 'relative'
      }}>
        <HighlightOffOutlined style={{ position: 'absolute', top: 5, right: 5, color: '#fff' }} onClick={() => onClose()}></HighlightOffOutlined>
        {images.length > 1 && (<>
          <KeyboardArrowLeft style={{
            fontSize: 45,
            position: 'absolute',
            top: '50%',
            color: '#fff'
          }} onClick={() => setPosition(position - 1 < 0 ? images.length - 1 : 0)}></KeyboardArrowLeft>
          <KeyboardArrowRight style={{
            fontSize: 45,
            position: 'absolute',
            top: '50%',
            right: 0,
            color: '#fff'
          }} onClick={() => setPosition(position + 1 >= images.length ? 0 : position + 1)}></KeyboardArrowRight>
          <div style={{
            position: 'absolute',
            bottom: 0,
            color: '#fff',
            display: 'flex',
            justifyContent: 'center',
            width: '100%'
          }}>
            <div style={{
              padding: 5,
              border: '1px solid #fff',
              borderRadius: 5,
              marginBottom: 5
            }}>{position + 1}/{images.length}</div>
          </div>
        </>)}
        {images[position] && <img src={images[position]} alt={`Message image ${position + 1}`} style={{ width: '100%', display: 'block' }} />}
      </div>
    </Dialog>
  );
}
