<template>
    <div class="background" @click="closeModal"></div>
    <div class="input-modal">
        <p class="modal-text">{{ message }}</p>
        <input v-model="input" type="text" ref="modalinput" @keydown.enter="submitInput" class="modal-input"
            placeholder="" />
        <button class="modal-button" @click="submitInput">OK</button>
        <button class="modal-button" @click="closeModal">Cancel</button>
    </div>
</template>

<script>
export default {
    data() {
        return {
            input: '',
        }
    },
    props: ['message'],
    methods: {
        closeModal() {
            this.$emit('show-modal', false);
        },
        submitInput() {
            const input = this.input.trim();
            if (input === '') {
                console.log('foo')
                return
            }
            console.log(input)
            this.$emit('input-content', input)
            this.closeModal()
        }
    },
    mounted() {
        this.$refs.modalinput.focus();
    }
}
</script>

<style scoped>
.background {
    position: fixed;
    right: 0px;
    top: 0px;
    height: 100%;
    width: 100%;
    background-color: black;
    z-index: 998;
    opacity: 0.2;

}

.modal-text {
    width: 95%;
    margin-left: auto;
    margin-right: auto;
    text-align: left;
}

.input-modal {
    border: none;
    border-radius: 4px;
    position: fixed;
    top: 50%;
    left: 50%;
    height: auto;
    width: 400px;
    color: #000000;
    transform: translate(-50%, -50%);
    background-color: #f2f2f2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    padding: 20px;
    z-index: 999;
}

.modal-input {
    width: 90%;
    display: block;
    margin-left: auto;
    margin-right: auto;
    padding: 10px;
    font-size: 16px;
    border-radius: 4px;
    border: none;
    background-color: #f2f2f2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.modal-button {
    float: right;
    position: relative;
    margin-top: 10px;
    margin-right: 10px;
    width: 20%;
    padding-top: 10px;
    padding-bottom: 10px;
    font-size: 16px;
    font-weight: bold;
    color: #ffffff;
    background-color: #169ba0;
    border: none;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    cursor: pointer;
}</style>