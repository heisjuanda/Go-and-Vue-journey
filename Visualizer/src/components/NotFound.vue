<script setup lang="ts">
import store from "../Store";

import notFoundIcon from "/Icons/notFound.png";
</script>

<template>
  <div class="flex flex-col mx-[auto] block w-[fit-content]">
    <img
      class="w-[100px] mx-[auto]"
      :src="notFoundIcon"
      alt="404 error when searching"
    />
    <p class="">The word: {{ store.searchedValue }} was not found</p>
  </div>
</template>
