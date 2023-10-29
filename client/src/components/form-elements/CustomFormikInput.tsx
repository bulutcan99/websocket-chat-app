import { Flex, FlexProps, theme } from "@chakra-ui/react";
import { Field } from "formik";
import { ReactNode } from "react";

interface Props {
  placeholder?: string;
  containerSettings?: FlexProps;
  inputStyle?: React.CSSProperties;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  type?: string;
  onChange?: (e: React.FormEvent<HTMLInputElement>) => void;
  onFocus?: (e: React.FormEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FormEvent<HTMLInputElement>) => void;
  id: string;
  name: string;
}

const CustomFormikInput = ({
  placeholder,
  containerSettings,
  inputStyle,
  leftIcon,
  onChange,
  onFocus,
  onBlur,
  rightIcon,
  type,
  id,
  name,
}: Props) => {
  return (
    <Flex
      alignItems={"center"}
      gap={"10px"}
      width={"100%"}
      paddingLeft={"10px"}
      paddingRight={"10px"}
      backgroundColor={"gray.100"}
      borderRadius={"2px"}
      padding={"10px"}
      {...containerSettings}
    >
      {leftIcon}
      <Field
        id={id}
        name={name}
        type={type}
        placeholder={placeholder}
        onChange={onChange}
        onFocus={onFocus}
        onBlur={onBlur}
        style={{
          border: "none !important",
          outline: "none !important",
          shadow: "none !important",
          backgroundColor: theme.colors.gray["100"],
          width: "100%",
          ...inputStyle,
        }}
      />
      {rightIcon}
    </Flex>
  );
};

export default CustomFormikInput;
