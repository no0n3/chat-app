import axios from "axios";
import { get, post } from "../utils/http";

export const getChats = (token: string) => get(`chat`, token);
export const getChat = (chatId: string, token: string) => get(`chat/${chatId}`, token);
export const getChatMessages = (chatId: string, token: string) => get(`chat/${chatId}/messages`, token);
export const getContacts = (token: string) => get('contacts', token);
export const getUser = (userId: string, token: string) => get(`user/${userId}`, token);
export const addContact = (userId: string, token: string) => post({
  path: `user/${userId}/add-contact`,
  token,
  payload: {}
});
export const removeContact = (userId: string, token: string) => post({
  path: `user/${userId}/remove-contact`,
  token,
  payload: {}
});
export const getProfiles = (token: string) => get('find', token);
export const uploadMedia = (formData: FormData, token: string) => post({
  path: 'upload',
  token,
  payload: formData,
  headers: {
    'Content-Type': 'multipart/form-data'
  }
});
export const uploadProfileImage = (formData: FormData, token: string) => post({
  path: 'upload-profile-image',
  token,
  payload: formData,
  headers: {
    'Content-Type': 'multipart/form-data'
  }
});
export const login = (payload: any) => axios.post(`${process.env.REACT_APP_ENDPOINT}/api/login`, payload);
export const signup = (payload: any) => axios.post(`${process.env.REACT_APP_ENDPOINT}/api/sign-up`, payload);
