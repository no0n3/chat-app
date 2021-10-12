import { Add, HighlightOffOutlined } from "@material-ui/icons";

export default function ImagePicker(props: any) {
  const {
    images,
    onAddHandler,
    onRemoveHandler
  } = props;

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'row'
    }}>
      {images.map((image: any) => (<div key={image.id} style={{
        padding: 1,
        height: 120,
        position: 'relative'
      }}>
        <img src={image.path} alt="Message image" style={{ height: '100%' }} />
        <HighlightOffOutlined
          onClick={() => onRemoveHandler(image)}
          style={{
            position: 'absolute',
            top: 5,
            right: 5
          }}
        ></HighlightOffOutlined>
      </div>))}
      <div style={{
        padding: 1,
        width: 120,
        height: 120,
        backgroundColor: 'gray',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
      }} onClick={() => onAddHandler()}>
        <Add style={{ fontSize: 60, color: '#fff' }}></Add>
      </div>
    </div>
  );
}
