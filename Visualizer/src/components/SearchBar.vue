<script setup lang="ts">
import { defineModel } from "vue";

import chevronIcon from "/Icons/chevron.png";
import searchIcon from "/Icons/search.png";
import store from "../Store";

const searchValue = defineModel<string | undefined | null>("searchTerm");
const SearchFunction = defineModel<() => Promise<void>>("SearchFunction");
const isLoading = defineModel<boolean>("isLoading");

const handleChangeOrder = () => {
  store.setOrder(!store.order);
  SearchFunction.value && SearchFunction.value();
};
</script>

<template>
  <div
    class="flex flex-row-reverse justify-center items-center gap-[15px] mt-[20px]"
  >
    <button
      :onclick="handleChangeOrder"
      class="text-white bg-[#535C91] hover:bg-blue-800 active:scale-110 transition-[0.3s] font-medium rounded-full text-sm p-2.5 text-center inline-flex items-center me-2 dark:bg-[#535C91] dark:hover:bg-blue-700"
    >
      <span class="ml-[2px] mr-[10px]">Date</span>
      <img
        :class="`w-[18px] transition-[0.25s] transform ${
          store.order ? 'rotate-180' : 'rotate-0'
        }`"
        :src="chevronIcon"
        alt="chevron for order"
      />
    </button>
    <input
      class="bg-[#535C91] color-[white] rounded-[25px] text-[16px] pt-[5px] pb-[5px] pl-[10px] pr-[10px] w-[50%] min-w-[250px] max-w-[600px] truncate overflow-hidden"
      type="text"
      placeholder="Search..."
      v-model="searchValue"
      @keydown.enter="SearchFunction"
      name="search"
      id="search"
      minlength="2"
      :disabled="isLoading"
    />
    <label
      for="search"
      class="grid place-items-center hover:scale-110 transition-[0.3s]"
    >
      <button v-on:click="SearchFunction">
        <img class="w-7 h-7" :src="searchIcon" alt="search button" />
      </button>
    </label>
  </div>
</template>
