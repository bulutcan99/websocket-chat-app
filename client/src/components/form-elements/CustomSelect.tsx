import { Box, Flex, Text, FlexProps } from "@chakra-ui/react";
import { ReactNode, Dispatch, SetStateAction, useState } from "react";
import AngleDownIcon from "../icons/AngleDownIcon";
import { CustomSelectOptions } from "@/utils/types/types";

interface Props {
  containerSettings?: FlexProps;
  inputStyle?: React.CSSProperties;
  leftIcon?: ReactNode;
  selectedOption: { value: string; label: string };
  setSelectedOption: Dispatch<SetStateAction<{ value: string; label: string }>>;
  isOptionOpen: boolean;
  setIsOptionOpen: Dispatch<SetStateAction<boolean>>;
  options: CustomSelectOptions[];
}

const CustomSelect = ({
  containerSettings,
  leftIcon,
  selectedOption,
  setSelectedOption,
  isOptionOpen,
  setIsOptionOpen,
  options,
}: Props) => {
  const optionHeight = 28;

  return (
    <Flex
      alignItems={"center"}
      gap={"10px"}
      width={"100%"}
      backgroundColor={"gray.100"}
      borderRadius={"2px"}
      padding={"10px"}
      position={"relative"}
      userSelect={"none"}
      outline={isOptionOpen ? "1px solid midnightblue" : "none"}
      transition={"outline 500ms"}
      {...containerSettings}
    >
      {leftIcon}
      <Box width={"100%"}>
        <Flex
          as={"div"}
          alignItems={"center"}
          justifyContent={"space-between"}
          padding={"0px 10px"}
          cursor={"pointer"}
          onClick={() => {
            setIsOptionOpen((prev) => !prev);
          }}
        >
          <Text>{selectedOption?.label}</Text>
          <AngleDownIcon
            className="angle-down-icon"
            width={20}
            height={20}
            fill="midnightblue"
            transform={isOptionOpen ? "rotate(180)" : "rotate(0)"}
          />
        </Flex>
        <Box
          width={"100%"}
          position={"absolute"}
          top={"44px"}
          left={0}
          height={isOptionOpen ? optionHeight * options.length - 1 : "0px"}
          backgroundColor={"white"}
          overflow={"hidden"}
          zIndex={10}
          outline={isOptionOpen ? "1px solid midnightblue" : "none"}
          transition={"height 500ms, outline 500ms"}
        >
          {options.map((option) => {
            return (
              <Flex
                flexDirection={"column"}
                justifyContent={"center"}
                height={`${optionHeight}px`}
                backgroundColor={"gray.100"}
                padding={"8px 12px"}
                key={option.value}
                color={"black"}
                cursor={"pointer"}
                _hover={{ backgroundColor: "midnightblue", color: "white" }}
                transition={"background-color 300ms, color 300ms"}
                onClick={() => {
                  setSelectedOption(option);
                  setIsOptionOpen(false);
                }}
              >
                {option.label}
              </Flex>
            );
          })}
        </Box>
      </Box>
    </Flex>
  );
};

export default CustomSelect;
