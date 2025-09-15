<template>
  <div>
    <h2 class="text-xl font-bold mb-3">Co-Owner Dashboard</h2>
    <p>
      Co-owner has nearly owner permissions. Use Owner Dash for full control.
    </p>
    <OwnerPartial />
  </div>
</template>

<script lang="ts">
import OwnerPartial from "./ownerDash.vue";
export default { components: { OwnerPartial } };
</script>
