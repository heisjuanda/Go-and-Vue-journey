<script setup lang="ts">
import { defineModel } from "vue";

import { Hit } from "../types/emailTypes";

const searchValue = defineModel<String>("searchTerm");
const emailContent = defineModel<Hit>("emailContent");

let emailId: string | undefined;
if (emailContent.value?._source?.message_id) {
  emailId = emailContent.value._source.message_id;
}
</script>

<template>
  <RouterLink
    :to="{
      name: 'email-content',
      params: { email_id: emailId },
    }"
    class="email-link"
  >
    <div class="envelope-wrapper">
      <div class="lid one"></div>
      <div class="lid two"></div>
      <div class="envelope"></div>
      <div class="letter">
        <p>{{ searchValue || "..." }}</p>
      </div>
    </div>
  </RouterLink>
</template>

<style>
.email-link {
  display: block;

  transform: scale(0.8);
}
.email-link:hover .envelope-wrapper {
  transform: translateY(-5px);
}
.envelope-wrapper {
  height: 200px;
  width: 300px;

  background-color: #3760c9;

  position: relative;
  z-index: 0;

  display: flex;
  justify-content: center;

  transition: 0.3s ease;
}
.lid {
  position: absolute;
  top: 0;
  left: 0;
  
  height: 100%;
  width: 100%;
  
  border-right: 150px solid transparent;
  border-bottom: 100px solid transparent;
  border-left: 150px solid transparent;

  transform-origin: top;
  
  transition: transform 0.25s linear;
}

.lid.one {
  border-top: 100px solid #658ced;

  transform: rotateX(0deg);
  
  z-index: 3;
  
  transition-delay: 0.25s;
}

.lid.two {
  border-top: 100px solid #3760c9;

  transform: rotateX(90deg);
  
  z-index: 1;
  
  transition-delay: 0.1s;
}

.envelope {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 3;

  height: 100%;
  width: 100%;
  
  border-top: 100px solid transparent;
  border-right: 150px solid #c4dff0;
  border-bottom: 100px solid #c4dff0;
  border-left: 150px solid #a4d4f2;
}

.letter {
  position: absolute;
  top: 0;
  z-index: 2;

  width: 80%;
  height: 80%;
  
  background-color: white;
  
  border-radius: 15px;
  
  transition: 0.25s;
}

.letter p {
  text-align: center;
  font-size: 30px;
  
  margin-top: 30px;
  
  color: #3b4049;

  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  width: 50%;
  
  margin-inline: auto;
}

.envelope-wrapper:hover .lid.one {
  transform: rotateX(90deg);
  transition-delay: 0s;
}

.envelope-wrapper:hover .lid.two {
  transform: rotateX(180deg);
  transition-delay: 0.25s;
}

.envelope-wrapper:hover .letter {
  transform: translateY(-50px);
  transition-delay: 0.25s;
}
</style>
