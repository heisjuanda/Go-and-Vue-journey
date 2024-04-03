<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";

import Loader from "../../components/Loader.vue";
import NotFound from "../../components/NotFound.vue";

import {
  getEmailFooterDetails,
  getParsedEmailId,
  getEmailLastHeaderDetails,
} from "../../helpers/emailDetails";

import getSingleEmail from "../../helpers/fetchSingleEmail";

import { Source } from "../../types/emailTypes";

import store from "../../Store";

window.scrollTo(0, 0);

const route = useRoute();
const parsedEmailId: string | null | undefined = getParsedEmailId(route);

const emailSheetRef = ref<HTMLElement | null>(null);
const isLoadingRef = ref<boolean>(true);
const isErrorRef = ref<boolean>(false);

const allEmailContent = ref<Source | null | undefined>();
const timestamp = ref<string | null | undefined>();
const date = ref<string | null | undefined>();
const sent = ref<string | null | undefined>();
const dateSub = ref<string | null | undefined>();
const subject = ref<string | null | undefined>();
const body = ref<string | null | undefined>();
const cc = ref<string | null | undefined>();
const xCc = ref<string | null | undefined>();

watch([allEmailContent, emailSheetRef], () => {
  const handleHighLightText = () => {
    if (!emailSheetRef.value || !store?.searchedValue || !allEmailContent.value)
      return;
    const allFields = emailSheetRef.value.querySelectorAll("p");
    if (!allFields) return;
    for (const field of allFields) {
      let emailContent = field.innerHTML;
      const regex = new RegExp(`(${store.searchedValue})`, "gi");
      if (regex.test(emailContent)) {
        emailContent = emailContent.replace(
          regex,
          `<mark class="highlight">${store.searchedValue}</mark>`
        );
      }
      field.innerHTML = emailContent;
    }
  };
  handleHighLightText();
});

onMounted(() => {
  const loadEmailProps = () => {
    if (!allEmailContent.value) {
      isErrorRef.value = true;
      return;
    }
    timestamp.value = allEmailContent.value?.["@timestamp"];
    date.value = allEmailContent.value?.date;
    sent.value = allEmailContent.value?.sent;
    dateSub.value = allEmailContent.value?.date_subemail;
    subject.value = allEmailContent.value?.subject;
    body.value = allEmailContent.value?.body;
    cc.value = allEmailContent.value?.cc;
    xCc.value = allEmailContent.value?.x_cc;
  };

  const getEmail = async () => {
    if (!parsedEmailId) return;
    isLoadingRef.value = true;
    allEmailContent.value = await getSingleEmail(parsedEmailId);
    loadEmailProps();
    isLoadingRef.value = false;
  };
  getEmail();
});
</script>

<template>
  <section
    v-if="!isLoadingRef && !isErrorRef"
    ref="emailSheetRef"
    class="email-sheet opacity-0 grid overflow-hidden pt-[20px] pl-[25px] pb-[20px] pr-[25px] mt-[5vh] mb-[5vh] mx-[auto] w-[80%] max-w-[80%] md:max-w-[650px] min-h-[120vh] bg-white"
  >
    <header
      class="flex flex-col justify-center align-center w-[80%] max-w-[80%] md:max-w-[650px]"
    >
      <div v-if="parsedEmailId" class="mb-[10px]">
        <p class="text-black text-[10px]">
          {{ parsedEmailId }}
        </p>
      </div>

      <div class="flex flex-col mt-[5px] mb-[5px]">
        <div v-if="timestamp">
          <p class="text-black text-[10px]">
            {{ new Date(timestamp) }}
          </p>
        </div>
        <div v-if="date || sent">
          <p class="text-black text-[10px]">{{ date || sent }}</p>
        </div>
        <div v-if="dateSub">
          <p class="text-black text-[10px]">{{ dateSub }}</p>
        </div>
      </div>
    </header>

    <header class="mt-[30px] mb-[20px] flex flex-col align-center">
      <div>
        <b>
          <p
            v-if="subject"
            class="text-center mb-[30px] text-black text-[16px]"
          >
            {{ subject }}
          </p>
        </b>
      </div>

      <div
        v-for="(lastHeaderDetails, _) in getEmailLastHeaderDetails(
          allEmailContent
        )"
      >
        <div
          v-if="
            lastHeaderDetails?.detailFirstValue ||
            lastHeaderDetails?.detailSecondValue
          "
          class="mb-[5px] flex flex-row flex-wrap gap-[4px] text-[12px]"
        >
          <b>
            <p class="text-black">{{ lastHeaderDetails.detailName }}:</p>
          </b>
          <p v-if="lastHeaderDetails.detailFirstValue" class="text-black">
            {{ lastHeaderDetails.detailFirstValue }}
          </p>
          <p v-if="lastHeaderDetails.detailSecondValue" class="text-black">
            {{ lastHeaderDetails.detailSecondValue }}
          </p>
        </div>
      </div>
    </header>

    <article class="email-body mt-[20px] max-w-[600px]]">
      <div v-if="body">
        <p class="text-black text-[14px]">
          {{ body }}
        </p>
      </div>
    </article>

    <footer class="mt-[20px]">
      <div v-if="cc || xCc" class="flex flex-col gap-[5px]">
        <b>
          <p class="text-black text-[12px]">cc:</p>
        </b>
        <p v-if="cc" class="text-black text-[12px]">
          {{ cc }}
        </p>
        <p v-if="xCc" class="text-black text-[12px]">
          {{ xCc }}
        </p>
      </div>

      <div
        class="mt-[10px] flex flex-col gap-[5px]"
        v-for="(details, _) in getEmailFooterDetails(allEmailContent)"
      >
        <div
          v-if="details?.detailValue"
          class="mb-[10px] flex flex-col gap-[5px]"
        >
          <div class="flex flex-row gap-[4px]">
            <b>
              <p class="text-black text-[12px]">{{ details.detailName }}:</p>
            </b>
            <p class="text-black text-[12px]">
              {{ details.detailValue }}
            </p>
          </div>
        </div>
      </div>
    </footer>
  </section>
  <div
    v-else-if="isLoadingRef && !isErrorRef"
    class="flex items-center justify-center h-[100dvh] w-[100vw]"
  >
    <Loader />
  </div>
  <div
    v-else-if="isErrorRef"
    class="flex items-center justify-center h-[100dvh] w-[100vw]"
  >
    <NotFound :isEmailError="isErrorRef" />
  </div>
</template>

<style>
@keyframes revealSheet {
  from {
    transform: translateY(100%) rotate(60deg);
    opacity: 0;
  }
  to {
    transform: translateY(0%) rotate(0deg);
    opacity: 1;
  }
}
@keyframes revealSheetMobile {
  from {
    transform: translateX(-100%);
    opacity: 0;
  }
  to {
    transform: translateX(0%);
    opacity: 1;
  }
}
.email-sheet {
  animation: revealSheet 0.7s ease-in-out 0.2s forwards;
}
.email-body {
  max-width: 600px !important;
  width: calc(80vw - 50px);
}
.highlight {
  background-color: orange;
  color: black;
}

@media screen and (max-width: 550px) {
  .email-sheet {
    animation: revealSheetMobile 0.7s ease-in-out 0.2s forwards;

    min-width: 95% !important;
  }
  .email-body {
    width: calc(95vw - 50px);
  }
}
</style>
