<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue";

import store from "../../Store";

import SearchBar from "../../components/SearchBar.vue";
import EmailPreview from "../../components/EmailPreview.vue";
import Loader from "../../components/Loader.vue";
import NotFound from "../../components/NotFound.vue";

import search from "../../helpers/fetchEmails";

import { Hits, Hit } from "../../types/emailTypes";

import { PAGINATION_VALUE } from "../../constants/constants";

import FindMail from "/Icons/FindMail.png";

const searchValue = ref<string>(store.searchedValue);
const isLoading = ref<boolean>(false);
const notFound = ref<boolean>(false);

const hasValidSearch = computed(
  () =>
    !isLoading.value &&
    Array.isArray(store.fetchedEmails) &&
    store.fetchedEmails.length > 0
);

const handleResetPagination = () => {
  store.setHasMoreResults(true);
  store.resetPagination();
};

const handleSearch = async () => {
  let isUserSearching = true;
  if (!searchValue.value || searchValue.value.length < 2) {
    store.setIsEmailResponseForPagination(true);
    isLoading.value = false;
    notFound.value = false;
    isUserSearching = false;
  }

  store.setIsSearching(true);
  isLoading.value = true;
  const oldSearchedTerm = store.searchedValue;
  store.setSearchedValue(searchValue.value);

  if (searchValue.value !== oldSearchedTerm) handleResetPagination();

  const hitsInformation: Hits | undefined = await search(
    searchValue.value,
    String(store.pagination),
    store.order ? "-" : "",
    isUserSearching
  );
  let allResults: Hit[] | undefined | null = [];
  if (hitsInformation?.hits) allResults = hitsInformation.hits;
  if (
    hitsInformation?.total?.value &&
    hitsInformation.total.value <= store.pagination
  ) {
    store.setHasMoreResults(false);
  }

  if (
    store.hasMoreResults &&
    hitsInformation?.total?.value &&
    allResults.length >= hitsInformation.total.value
  ) {
    store.setIsEmailResponseForPagination(false);
  }

  if (!allResults || allResults?.length < 1) {
    store.setFetchedEmails([]);
    handleResetPagination();
    store.setIsEmailResponseForPagination(true);
    isLoading.value = false;
    notFound.value = true;
    return;
  }

  if (
    !isUserSearching &&
    Array.isArray(store.fetchedEmails) &&
    allResults.length > 0
  ) {
    store.setFetchedEmails(allResults);
    if (hitsInformation?.total?.value)
      store.setIsEmailResponseForPagination(
        allResults.length < hitsInformation?.total?.value
      );
    isLoading.value = false;
    return;
  }

  if (Array.isArray(allResults) && allResults?.length > 0) {
    store.setFetchedEmails(allResults);
    if (hitsInformation?.total?.value)
      store.setIsEmailResponseForPagination(
        allResults.length < hitsInformation?.total?.value
      );
  }

  isLoading.value = false;
};

const handlePagination = () => {
  store.setPagination(PAGINATION_VALUE);
  handleSearch();
};

watch(searchValue, (element) => {
  if (!element) {
    store.setIsSearching(false);
    handleResetPagination();
    notFound.value = false;
  }
});

onMounted(() => {
  if (Array.isArray(store.fetchedEmails) && store.fetchedEmails.length < 1)
    handleSearch();
});
</script>

<template>
  <section
    ref="allEmailsView"
    class="flex flex-col gap-[20px] align-center justify-center"
  >
    <header class="w-[fit-content]">
      <h1
        class="text-[20px] flex flex-row-reverse gap-[5px] align-center justify-center ml-[20px] mt-[20px]"
      >
        <span>FindMail</span>
        <img
          class="h-[auto] max-h-[30px]"
          :src="FindMail"
          alt="FindMail logo"
        />
      </h1>
    </header>
    <SearchBar
      v-model:searchTerm="searchValue"
      v-model:SearchFunction="handleSearch"
      v-model:isLoading="isLoading"
    />
    <article class="mt-[100px] ml-[40px] mr-[40px]">
      <div
        v-if="isLoading"
        class="w-[50vw] h-[50vh] flex align-center justify-center mx-[auto]"
      >
        <Loader />
      </div>
      <ul
        v-if="hasValidSearch"
        class="align-center justify-center flex flex-row gap-[20px] flex-wrap"
      >
        <li v-for="(emailContent, index) in store.fetchedEmails" :key="index">
          <EmailPreview
            v-model:searchTerm="searchValue"
            :emailContent="emailContent?._source?.message_id"
          />
        </li>
      </ul>
      <div v-else-if="notFound">
        <NotFound />
      </div>

      <button
        class="px-3 py-2 text-xs font-medium text-center text-white bg-blue-700 rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 block w-fit-content mt-[40px] mb-[50px] mx-[auto] hover:scale-110 transition-[0.3s]"
        :onclick="handlePagination"
        :disabled="isLoading"
        v-if="
          hasValidSearch &&
          store.hasMoreResults &&
          store.isEmailResponseForPagination
        "
      >
        More Results
      </button>
    </article>
  </section>
</template>
