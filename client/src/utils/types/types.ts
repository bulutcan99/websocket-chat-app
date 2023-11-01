export type RegisterRequestBody = {
  nickname: string;
  name: string;
  surname: string;
  email: string;
  password: string;
  user_role?: string;
};

export type CustomSelectOptions = {
  value: number | string;
  label: string;
};
