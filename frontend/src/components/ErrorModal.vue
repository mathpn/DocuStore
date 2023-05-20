<template>
    <div class="error-modal">
        <p><b>Error: {{ errorMsg }}</b></p>
        <div id="slide-bar"></div>
    </div>
</template>

<script>
export default {
    props: ['errorMsg']
}
</script>

<style scoped>
.error-modal {
    border: none;
    border-radius: 4px;
    position: fixed;
    top: 50%;
    left: 50%;
    color: #521414;
    transform: translate(-50%, -50%);
    background-color: #ffa29f;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    padding: 20px;
    z-index: 999;
}

#slide-bar {
    border: none;
    border-radius: 4px;
    position: relative;
    height: 5px;
    background-color: #521414;
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    z-index: 998;
    animation: slideOut 1s linear forwards;
}

@keyframes slideOut {
    from {
        width: 100%;
    }

    to {
        width: 0;
    }
}
</style>