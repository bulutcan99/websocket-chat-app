import { RegisterRequestBody } from "./types/types";
import { urls } from "./urls";

export const fetchLogin = async () => {
  return await fetch(urls.login)
    .then((response) => {
      return response?.json();
    })
    .then((data) => {
      return data;
    })
    .catch((error) => {
      return error;
    });
};

export const fetchRegister = async (body: RegisterRequestBody) => {
  return await fetch(urls.register, {
    method: "POST",
    body: JSON.stringify({ ...body }),
  })
    .then((response) => {
      return response?.json();
    })
    .then((data) => {
      return data;
    })
    .catch((error) => {
      return error;
    });
};
