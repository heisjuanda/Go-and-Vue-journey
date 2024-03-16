import { Root } from "../types/emailTypes";

import {
  AUTHORIZATION_LEVEL,
  ENDPOINT,
  METHOD_GET,
  SERVER,
  USER,
  PASSWORD,
} from "../constants/constants";

const search = async (searchTerm: string, page: string, orderDate: string) => {
  const url = new URL(`${SERVER}${ENDPOINT}`);
  url.searchParams.append("term", searchTerm);
  url.searchParams.append("from", page);
  url.searchParams.append("order", orderDate);

  try {
    const response = await fetch(url, {
      method: METHOD_GET,
      headers: {
        Authorization: AUTHORIZATION_LEVEL + btoa(`${USER}:${PASSWORD}`),
      },
    });
    if (!response.ok) {
      throw new Error(
        `Error in your request to Get Emails, status: ${response?.status}`
      );
    }
    const data: Root = await response.json();
    if (data) return data?.hits;
  } catch (error) {
    console.error("There was a problem with the fetch operation:", error);
  }
};

export default search;
