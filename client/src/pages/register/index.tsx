import { Box, Text, Flex, Button } from "@chakra-ui/react";
import EmailIncon from "@/components/icons/EmailIncon";
import LockIcon from "@/components/icons/LockIcon";
import { FormEvent, useEffect, useState } from "react";
import EyeIcon from "@/components/icons/EyeIcon";
import EyeLineIcon from "@/components/icons/EyeLineIcon";
import { AllPages, InputType } from "@/utils/types/enums";
import Link from "next/link";
import Layout from "@/components/Layout";
import UserIcon from "@/components/icons/UserIcon";
import { fetchRegister } from "@/utils/api";
import { Form, Formik, FormikHelpers } from "formik";
import CustomFormikInput from "@/components/form-elements/CustomFormikInput";
import CustomSelect from "@/components/form-elements/CustomSelect";

export interface Values {
  nickname: string;
  name: string;
}

const Register = () => {
  const [emailValue, setEmailValue] = useState<string>("");
  const [passwordValue, setpasswordValue] = useState<string>("");
  const [nameSurnameValue, setNameSurnameValue] = useState<string>("");
  const [validation, setValidation] = useState({
    isMinLength: false,
    isUpperCharacter: false,
    isLowerCharacter: false,
    isNumber: false,
  });
  const [focusedInput, setFocusedInput] = useState({
    email: false,
    password: false,
    name: false,
    surname: false,
    nickname: false,
  });
  const [passwordInputType, setPasswordInputType] = useState(
    InputType.PASSWORD
  );
  const [selectedUserRole, setSelectedUserRole] = useState<{
    value: string;
    label: string;
  }>({ value: "user", label: "User" });
  const [isOptionOpen, setIsOptionOpen] = useState(false);

  useEffect(() => {
    const validationSchema = {
      isMinLength: false,
      isUpperCharacter: false,
      isLowerCharacter: false,
      isNumber: false,
    };
    if (passwordValue?.length < 9) {
      validationSchema.isMinLength = true;
    }
    if (/[A-Z]/.test(passwordValue)) {
      validationSchema.isUpperCharacter = true;
    }
    if (/[a-z]/.test(passwordValue)) {
      validationSchema.isLowerCharacter = true;
    }
    if (/\d/.test(passwordValue)) {
      validationSchema.isNumber = true;
    }
    setValidation(validationSchema);
  }, [passwordValue, passwordValue?.length]);

  const handleRegisterClick = () => {
    fetchRegister({
      name_surname: nameSurnameValue,
      email: emailValue,
      password: passwordValue,
      user_role: "string",
    });
  };

  const onSubmit = async (values: Values) => {
    const allValues = { ...values, user_role: selectedUserRole.value };
    console.log(allValues);
  };

  return (
    <Layout page={AllPages.REGISTER}>
      <Flex
        width={"100vw"}
        height={"100vh"}
        justifyContent={"center"}
        alignItems={"center"}
        backgroundColor={"cornflowerblue"}
      >
        <Box
          width={500}
          height={750}
          backgroundColor={"white"}
          borderRadius={6}
          boxShadow={"0px 0px 50px 1px midnightblue"}
          padding={"10px 80px"}
          animation={"startOpacityAnimation 1s linear"}
        >
          <Text
            display={"flex"}
            justifyContent={"center"}
            alignItems={"center"}
            width={"100%"}
            height={"20%"}
            fontSize={30}
            fontWeight={"bold"}
            letterSpacing={4}
            color={"midnightblue"}
            animation={"opacityAnimation 3s linear infinite"}
            userSelect={"none"}
          >
            REGISTER
          </Text>
          <Formik
            style={{ height: "60%" }}
            initialValues={{
              nickname: "",
              name: "",
              surname: "",
              email: "",
              password: "",
            }}
            onSubmit={onSubmit}
          >
            {({ values, handleChange, setFieldValue }) => (
              <Form>
                <CustomFormikInput
                  name="nickname"
                  id="nickname"
                  placeholder={"Nickname"}
                  value={values.nickname}
                  leftIcon={
                    <Box
                      opacity={focusedInput?.nickname ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <UserIcon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  type={"text"}
                  onFocus={() => {
                    setFocusedInput((prev) => ({ ...prev, nickname: true }));
                  }}
                  onBlur={() => {
                    setFocusedInput((prev) => ({
                      ...prev,
                      nickname: false,
                    }));
                  }}
                  onChange={handleChange}
                  containerSettings={{
                    outline: focusedInput?.nickname
                      ? "1px solid midnightblue"
                      : "none",
                    transition: "outline 500ms",
                  }}
                />
                <CustomFormikInput
                  id="name"
                  name="name"
                  placeholder={"Name"}
                  value={values.name}
                  leftIcon={
                    <Box
                      opacity={focusedInput?.name ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <UserIcon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  type={"text"}
                  onChange={handleChange}
                  onFocus={() => {
                    setFocusedInput((prev) => ({ ...prev, name: true }));
                  }}
                  onBlur={() => {
                    setFocusedInput((prev) => ({
                      ...prev,
                      name: false,
                    }));
                  }}
                  containerSettings={{
                    marginTop: "16px",
                    outline: focusedInput?.name
                      ? "1px solid midnightblue"
                      : "none",
                    transition: "outline 500ms",
                  }}
                />
                <CustomFormikInput
                  name="surname"
                  id="surname"
                  placeholder={"Surname"}
                  leftIcon={
                    <Box
                      opacity={focusedInput?.surname ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <UserIcon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  value={values.surname}
                  type={"text"}
                  onChange={handleChange}
                  onFocus={() => {
                    setFocusedInput((prev) => ({ ...prev, surname: true }));
                  }}
                  onBlur={() => {
                    setFocusedInput((prev) => ({
                      ...prev,
                      surname: false,
                    }));
                  }}
                  containerSettings={{
                    marginTop: "16px",
                    outline: focusedInput?.surname
                      ? "1px solid midnightblue"
                      : "none",
                    transition: "outline 500ms",
                  }}
                />
                <CustomFormikInput
                  name="email"
                  id="email"
                  placeholder={"Email"}
                  value={values.email}
                  leftIcon={
                    <Box
                      opacity={focusedInput?.email ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <EmailIncon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  type={"text"}
                  onChange={handleChange}
                  onFocus={() => {
                    setFocusedInput((prev) => ({ ...prev, email: true }));
                  }}
                  onBlur={() => {
                    setFocusedInput((prev) => ({ ...prev, email: false }));
                  }}
                  containerSettings={{
                    marginTop: "16px",
                    outline: focusedInput?.email
                      ? "1px solid midnightblue"
                      : "none",
                    transition: "outline 500ms",
                  }}
                />
                <CustomFormikInput
                  name="password"
                  id="password"
                  placeholder={"Password"}
                  value={values.password}
                  leftIcon={
                    <Box
                      opacity={focusedInput?.password ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <LockIcon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  rightIcon={
                    passwordInputType === InputType.PASSWORD ? (
                      <Box
                        userSelect={"none"}
                        cursor={"pointer"}
                        onClick={() => {
                          setPasswordInputType(InputType.TEXT);
                        }}
                      >
                        <EyeIcon
                          width={20}
                          height={20}
                          fill="midnightblue"
                          opacity={0.8}
                        />
                      </Box>
                    ) : (
                      <Box
                        userSelect={"none"}
                        cursor={"pointer"}
                        onClick={() => {
                          setPasswordInputType(InputType.PASSWORD);
                        }}
                      >
                        <EyeLineIcon
                          width={20}
                          height={20}
                          fill="midnightblue"
                          opacity={0.8}
                        />
                      </Box>
                    )
                  }
                  type={passwordInputType}
                  onChange={(e: FormEvent<HTMLInputElement>) => {
                    setpasswordValue(e.currentTarget.value);
                    setFieldValue("password", e.currentTarget.value);
                  }}
                  onFocus={() => {
                    setFocusedInput((prev) => ({ ...prev, password: true }));
                  }}
                  onBlur={() => {
                    setFocusedInput((prev) => ({ ...prev, password: false }));
                  }}
                  containerSettings={{
                    marginTop: "16px",
                    outline: focusedInput?.password
                      ? "1px solid midnightblue"
                      : "none",
                    transition: "outline 500ms",
                  }}
                />
                <Flex
                  flexDirection={"column"}
                  alignItems={"flex-start"}
                  justifyContent={"center"}
                  overflow={"hidden"}
                  height={passwordValue?.length ? 95 : 0}
                  transition={"height 500ms"}
                  userSelect={"none"}
                  marginTop={"6px"}
                >
                  <Text color={validation?.isMinLength ? "red" : "green"}>
                    Must be at least 9 characters long.
                  </Text>
                  <Text color={validation?.isUpperCharacter ? "green" : "red"}>
                    Must contain at least one uppercase letter.
                  </Text>
                  <Text color={validation?.isLowerCharacter ? "green" : "red"}>
                    Must contain at least one lowercase letter.
                  </Text>
                  <Text color={validation?.isNumber ? "green" : "red"}>
                    Must contain at least one digit.
                  </Text>
                </Flex>
                <CustomSelect
                  selectedOption={selectedUserRole}
                  setSelectedOption={setSelectedUserRole}
                  isOptionOpen={isOptionOpen}
                  setIsOptionOpen={setIsOptionOpen}
                  options={[
                    { value: "user", label: "User" },
                    { value: "admin", label: "Admin" },
                  ]}
                  leftIcon={
                    <Box
                      opacity={isOptionOpen ? 0.8 : 0.5}
                      transition={"opacity 500ms"}
                    >
                      <LockIcon width={20} height={20} fill="midnightblue" />
                    </Box>
                  }
                  containerSettings={{ marginTop: "16px" }}
                />
                <Button
                  width={"100%"}
                  color={"black"}
                  type="submit"
                  marginTop={"26px"}
                  backgroundColor={"gray.100"}
                  transition={"color 500ms, background-color 500ms"}
                  _hover={{ backgroundColor: "midnightblue", color: "white" }}
                  onClick={handleRegisterClick}
                >
                  REGISTER NOW
                </Button>
              </Form>
            )}
          </Formik>
          <Flex
            justifyContent={"center"}
            alignItems={"flex-end"}
            height={"20%"}
          >
            Already a member! Let's
            <Text
              textDecoration={"underline"}
              marginLeft={1}
              fontWeight={"bold"}
              cursor={"pointer"}
            >
              <Link href={"/"}>login.</Link>
            </Text>
          </Flex>
        </Box>
      </Flex>
    </Layout>
  );
};

export default Register;
