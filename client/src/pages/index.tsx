import Layout from "@/components/Layout";
import CustomFormikInput from "@/components/form-elements/CustomFormikInput";
import EmailIncon from "@/components/icons/EmailIncon";
import EyeIcon from "@/components/icons/EyeIcon";
import EyeLineIcon from "@/components/icons/EyeLineIcon";
import LockIcon from "@/components/icons/LockIcon";
import { AllPages, InputType } from "@/utils/types/enums";
import { Box, Button, Flex, Link, Text } from "@chakra-ui/react";
import { Form, Formik } from "formik";
import Head from "next/head";
import { useCallback, useState } from "react";

export default function Home() {
  const [passwordInputType, setPasswordInputType] = useState(
    InputType.PASSWORD
  );
  const [focusedInput, setFocusedInput] = useState({
    email: false,
    password: false,
    name: false,
    surname: false,
    nickname: false,
  });
  const initialValues = { email: "", password: "" };

  const onSubmit = useCallback(
    (values: { email: string; password: string }) => {
      console.log(values);
    },
    []
  );

  return (
    <>
      <Head>
        <title>Create Next App</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Layout page={AllPages.LOGIN}>
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
              height={"35%"}
              fontSize={30}
              fontWeight={"bold"}
              letterSpacing={4}
              color={"midnightblue"}
              animation={"opacityAnimation 3s linear infinite"}
              userSelect={"none"}
            >
              LOGIN
            </Text>
            <Formik initialValues={initialValues} onSubmit={onSubmit}>
              {({ values, handleChange, setFieldValue }) => (
                <Form>
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
                        <EmailIncon
                          width={20}
                          height={20}
                          fill="midnightblue"
                        />
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
                    onChange={handleChange}
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

                  <Button
                    width={"100%"}
                    color={"black"}
                    type="submit"
                    marginTop={"26px"}
                    backgroundColor={"gray.100"}
                    transition={"color 500ms, background-color 500ms"}
                    _hover={{ backgroundColor: "midnightblue", color: "white" }}
                  >
                    LOGIN
                  </Button>
                </Form>
              )}
            </Formik>
            <Flex
              justifyContent={"center"}
              alignItems={"flex-end"}
              height={"40%"}
            >
              Still not a member? Let's
              <Text
                textDecoration={"underline"}
                marginLeft={1}
                fontWeight={"bold"}
                cursor={"pointer"}
                color={"black"}
                _hover={{ color: "midnightblue" }}
                transition={"color 500ms"}
              >
                <Link href={"/register"}>Register!</Link>
              </Text>
            </Flex>
          </Box>
        </Flex>
      </Layout>
    </>
  );
}
