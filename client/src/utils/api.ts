import { RegisterRequestBody } from "./types/types";
import { urls } from "./urls";

export const fetchLogin = async () => {
  try {
    const response = await fetch(urls.login);
    const data = await response.json();
    return data;
  } catch (error) {
    return error;
  }
};

export const fetchRegister = async (body: RegisterRequestBody) => {
  try {
    const response = await fetch(urls.register, {
      method: "POST",
      body: JSON.stringify({ ...body }),
    });
    const data = await response.json();
    return data;
  } catch (error) {
    return error;
  }
};
