import axios from "axios";

const BASE_URL = `${process.env.REACT_APP_ENDPOINT}/api/`;

export function get(path: string, token: string) {
  return axios.get(BASE_URL + path, {
    headers: {
      'x-auth-token': token
    }
  })
    .then(response => response.data);;
}

export function post({
  path,
  token,
  payload,
  headers
}: {
  path: string,
  token: string,
  payload: any,
  headers?: any
}) {
  if (!payload) {
    payload = {};
  }

  if (!headers) {
    headers = {};
  }

  return axios.post(BASE_URL + path, payload, {
    headers: {
      'x-auth-token': token,
      ...headers
    }
  })
    .then(response => response.data);
}
