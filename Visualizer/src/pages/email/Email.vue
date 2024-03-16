<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";

import {
  getEmailFooterDetails,
  getParsedEmail,
  getEmailLastHeaderDetails,
} from "../../helpers/emailDetails";

import store from "../../Store";

import { Source } from "../../types/emailTypes";

window.scrollTo(0, 0);

const emailSheetRef = ref<HTMLElement | null>(null);

const route = useRoute();

const parsedEmail: Source | null | undefined = getParsedEmail(route);

const timestamp = parsedEmail?.["@timestamp"];
const id = parsedEmail?.message_id;

const date = parsedEmail?.date;
const sent = parsedEmail?.sent;
const dateSub = parsedEmail?.date_subemail;

const subject = parsedEmail?.subject;

const body = parsedEmail?.body;

const cc = parsedEmail?.cc;
const xCc = parsedEmail?.x_cc;
onMounted(() => {
  if (!emailSheetRef.value || !store?.searchedValue) return;
  let emailContent = emailSheetRef.value.innerHTML;
  const regex = new RegExp(`(${store.searchedValue})`, "gi");
  if (regex.test(emailContent)) {
    emailContent = emailContent.replace(
      regex,
      `<mark class="highlight">${store.searchedValue}</mark>`
    );
  }
  emailSheetRef.value.innerHTML = emailContent;
});
</script>

<template>
  <section
    ref="emailSheetRef"
    class="email-sheet opacity-0 grid overflow-hidden pt-[20px] pl-[25px] pb-[20px] pr-[25px] mt-[5vh] mb-[5vh] mx-[auto] w-[80%] max-w-[80%] md:max-w-[650px] min-h-[120vh] bg-white"
  >
    <header
      class="flex flex-col justify-center align-center w-[80%] max-w-[80%] md:max-w-[650px]"
    >
      <div v-if="id" class="mb-[10px]">
        <p class="text-black text-[10px]">
          {{ id }}
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
        <h3 v-if="subject" class="text-center mb-[30px] text-black text-[16px]">
          {{ subject }}
        </h3>
      </div>

      <div
        v-for="(lastHeaderDetails, _) in getEmailLastHeaderDetails(parsedEmail)"
      >
        <div
          v-if="
            lastHeaderDetails?.detailFirstValue ||
            lastHeaderDetails?.detailSecondValue
          "
          class="mb-[5px] flex flex-row flex-wrap gap-[4px] text-[12px]"
        >
          <p class="text-black">
            <b>{{ lastHeaderDetails.detailName }}:</b>
          </p>
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
        <p class="text-black text-[12px]">
          <b>cc:</b>
        </p>
        <p v-if="cc" class="text-black text-[12px]">
          {{ cc }}
        </p>
        <p v-if="xCc" class="text-black text-[12px]">
          {{ xCc }}
        </p>
      </div>

      <div
        class="mt-[10px] flex flex-col gap-[5px]"
        v-for="(details, _) in getEmailFooterDetails(parsedEmail)"
      >
        <div
          v-if="details?.detailValue"
          class="mb-[10px] flex flex-col gap-[5px]"
        >
          <div class="flex flex-row gap-[4px]">
            <p class="text-black text-[12px]">
              <b>{{ details.detailName }}:</b>
            </p>
            <p class="text-black text-[12px]">
              {{ details.detailValue }}
            </p>
          </div>
        </div>
      </div>
    </footer>
  </section>
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
