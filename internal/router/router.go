package router

import (
	"crypto/subtle"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const tmplSource string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8" />

    <title>Skpr Environment</title>

    <style>
        html {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
            background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg width='482' height='482' xmlns='http://www.w3.org/2000/svg' fill='%23E6E6E6'%3E%3Cpath d='M255.694 466.88a2.93 2.93 0 11-5.851.291 2.93 2.93 0 015.85-.292zm-25.56-2.784l.173.004a2.93 2.93 0 11-.342.002l.17-.006zm48.025.242a2.93 2.93 0 11-5.792.873 2.93 2.93 0 015.792-.873zm-70.056-2.46a2.93 2.93 0 11-.874 5.793 2.93 2.93 0 01.874-5.793zm92.158-2.306a2.93 2.93 0 11-5.677 1.445 2.93 2.93 0 015.677-1.445zm-114.031-2.116a2.93 2.93 0 11-1.446 5.677 2.93 2.93 0 011.446-5.677zm135.547-4.827a2.93 2.93 0 11-5.504 2.004 2.93 2.93 0 015.504-2.004zm-156.871-1.75a2.93 2.93 0 11-2.004 5.505 2.93 2.93 0 012.004-5.505zm177.59-7.3a2.93 2.93 0 11-5.278 2.541 2.93 2.93 0 015.278-2.541zm-198.153-1.368a2.93 2.93 0 11-2.542 5.278 2.93 2.93 0 012.542-5.278zm111.138 1.956a2.778 2.778 0 11-5.547.306 2.778 2.778 0 015.547-.306zm-25.276-2.625l.17.004a2.778 2.778 0 11-.337.002l.167-.006zm47.593-.153a2.778 2.778 0 11-5.48.914 2.778 2.778 0 015.48-.914zm-69.469-2.283a2.778 2.778 0 11-.915 5.48 2.778 2.778 0 01.915-5.48zm-21.643-4.85a2.778 2.778 0 11-1.512 5.346 2.778 2.778 0 011.512-5.346zm112.987 1.917a2.778 2.778 0 11-5.347 1.512 2.778 2.778 0 015.346-1.512zm62.538-3.663a2.93 2.93 0 11-5 3.055 2.93 2.93 0 015-3.055zm-237.466-.972a2.93 2.93 0 11-3.054 4.999 2.93 2.93 0 013.054-4.999zm196.096-2.956a2.778 2.778 0 11-5.147 2.091 2.778 2.778 0 015.147-2.09zm-155.134-1.528a2.778 2.778 0 11-2.09 5.148 2.778 2.778 0 012.09-5.148zm75.758-7.728a2.627 2.627 0 110 5.253 2.627 2.627 0 010-5.253zm139.26.208a2.93 2.93 0 11-4.67 3.537 2.93 2.93 0 014.67-3.537zm-274.418-.566a2.93 2.93 0 11-3.537 4.67 2.93 2.93 0 013.538-4.67zm160.501 1.225a2.627 2.627 0 11-5.212.659 2.627 2.627 0 015.212-.659zm-47.752-2.277a2.627 2.627 0 11-.658 5.212 2.627 2.627 0 01.658-5.212zm-73.407-.328a2.778 2.778 0 11-2.645 4.886 2.778 2.778 0 012.645-4.886zm195.397 1.121a2.778 2.778 0 11-4.886 2.644 2.778 2.778 0 014.886-2.644zm-51.92-3.11a2.627 2.627 0 11-5.09 1.307 2.627 2.627 0 015.09-1.306zM197 413.71a2.627 2.627 0 11-1.306 5.088A2.627 2.627 0 01197 413.71zm113.69-5.463a2.627 2.627 0 11-4.884 1.934 2.627 2.627 0 014.884-1.934zm-135.044-1.475a2.627 2.627 0 11-1.934 4.884 2.627 2.627 0 011.934-4.884zm-48.894-.798a2.778 2.778 0 11-3.166 4.566 2.778 2.778 0 013.166-4.566zm233.29.7a2.778 2.778 0 11-4.567 3.166 2.778 2.778 0 014.566-3.166zm37.813-1.894a2.93 2.93 0 11-4.295 3.984 2.93 2.93 0 014.295-3.984zm-308.643-.155a2.93 2.93 0 11-3.985 4.294 2.93 2.93 0 013.985-4.294zm66.119-7.413a2.627 2.627 0 11-2.531 4.604 2.627 2.627 0 012.53-4.604zm175.833 1.036a2.627 2.627 0 11-4.604 2.531 2.627 2.627 0 014.604-2.53zm-89.7-1.447a2.476 2.476 0 110 4.951 2.476 2.476 0 010-4.951zm-22.24-1.59a2.476 2.476 0 11-.704 4.9 2.476 2.476 0 01.705-4.9zm47.281 2.098a2.475 2.475 0 11-4.9.704 2.475 2.475 0 014.9-.704zM109.25 392.35a2.778 2.778 0 11-3.65 4.19 2.778 2.778 0 013.65-4.19zm268.349.27a2.778 2.778 0 11-4.19 3.649 2.778 2.778 0 014.19-3.649zm-180.16-2.15a2.476 2.476 0 11-1.394 4.75 2.476 2.476 0 011.395-4.75zm91.123 1.678a2.476 2.476 0 11-4.75 1.395 2.476 2.476 0 014.75-1.395zm124.87-3.757a2.93 2.93 0 11-3.877 4.392 2.93 2.93 0 013.877-4.392zm-339.8.258a2.93 2.93 0 11-4.393 3.876 2.93 2.93 0 014.393-3.876zm62.742-3.469a2.627 2.627 0 11-3.088 4.25 2.627 2.627 0 013.088-4.25zm213.85.581a2.627 2.627 0 11-4.25 3.088 2.627 2.627 0 014.25-3.088zM176.55 382.68a2.476 2.476 0 11-2.057 4.503 2.476 2.476 0 012.057-4.503zm133.11 1.223a2.475 2.475 0 11-4.504 2.057 2.475 2.475 0 014.503-2.057zm-216.31-7.02a2.778 2.778 0 11-4.087 3.763 2.778 2.778 0 014.087-3.763zm300.152-.162a2.778 2.778 0 11-3.762 4.088 2.778 2.778 0 013.762-4.088zm-152.038-2.447a2.324 2.324 0 110 4.649 2.324 2.324 0 010-4.649zm24.687.086a2.324 2.324 0 11-4.585.765 2.324 2.324 0 014.585-.765zm-46.7-1.91a2.324 2.324 0 11-.765 4.585 2.324 2.324 0 01.765-4.585zm-62.47-.456a2.476 2.476 0 11-2.677 4.165 2.476 2.476 0 012.677-4.165zm172.387.744a2.476 2.476 0 11-4.166 2.677 2.476 2.476 0 014.166-2.677zm97.93-2.204a2.93 2.93 0 11-3.42 4.756 2.93 2.93 0 013.42-4.756zm-367.58.668a2.93 2.93 0 11-4.757 3.42 2.93 2.93 0 014.756-3.42zm59.355-.334a2.627 2.627 0 11-3.596 3.83 2.627 2.627 0 013.596-3.83zm248.494.117a2.627 2.627 0 11-3.83 3.596 2.627 2.627 0 013.83-3.596zm-79.725-2.514a2.324 2.324 0 11-4.397 1.51 2.324 2.324 0 014.397-1.51zm-89.803-1.443a2.324 2.324 0 11-1.51 4.397 2.324 2.324 0 011.51-4.397zm-118.79-7.266a2.778 2.778 0 11-4.476 3.29 2.778 2.778 0 014.477-3.29zm328.312-.593a2.778 2.778 0 11-3.29 4.477 2.778 2.778 0 013.29-4.477zm-268.429-.537a2.476 2.476 0 11-3.242 3.742 2.476 2.476 0 013.242-3.742zm208.155.25a2.476 2.476 0 11-3.742 3.242 2.476 2.476 0 013.742-3.242zm-169.476-.727a2.324 2.324 0 11-2.213 4.089 2.324 2.324 0 012.213-4.089zm130.456.938a2.324 2.324 0 11-4.089 2.212 2.324 2.324 0 014.089-2.212zm74.655-4.942a2.627 2.627 0 11-3.349 4.048 2.627 2.627 0 013.35-4.048zm-279.219.35a2.627 2.627 0 11-4.048 3.348 2.627 2.627 0 014.048-3.348zm335.615-3.116a2.93 2.93 0 11-2.929 5.073 2.93 2.93 0 012.93-5.073zM47.61 352.457a2.93 2.93 0 11-5.073 2.929 2.93 2.93 0 015.073-2.93zm207.487.663a2.173 2.173 0 11-4.324.44 2.173 2.173 0 014.324-.44zm-25.04-1.953l.155.01a2.173 2.173 0 11-.308-.01h.153zm-21.886-4.52a2.173 2.173 0 11-1.302 4.148 2.173 2.173 0 011.302-4.147zm69.31 1.424a2.173 2.173 0 11-4.148 1.3 2.173 2.173 0 014.148-1.3zm-118.161-1.997a2.324 2.324 0 11-2.855 3.668 2.324 2.324 0 012.855-3.668zm167.55.407a2.324 2.324 0 11-3.669 2.855 2.324 2.324 0 013.668-2.855zm36.183-3.864a2.476 2.476 0 11-3.243 3.742 2.476 2.476 0 013.243-3.742zm-239.686.25a2.476 2.476 0 11-3.742 3.242 2.476 2.476 0 013.742-3.242zm-56.249-1.673a2.778 2.778 0 11-4.811 2.778 2.778 2.778 0 014.811-2.778zm352.486-1.017a2.778 2.778 0 11-2.779 4.812 2.778 2.778 0 012.779-4.812zm-232.112-2.403a2.173 2.173 0 11-2.11 3.8 2.173 2.173 0 012.11-3.8zm110.898.845a2.173 2.173 0 11-3.8 2.11 2.173 2.173 0 013.8-2.11zm97.657-3.093a2.627 2.627 0 11-2.815 4.435 2.627 2.627 0 012.815-4.435zm-305.541.81a2.627 2.627 0 11-4.436 2.815 2.627 2.627 0 014.436-2.815zm358.866-5.202a2.93 2.93 0 11-2.41 5.34 2.93 2.93 0 012.41-5.34zM37.43 332.599a2.93 2.93 0 11-5.34 2.41 2.93 2.93 0 015.34-2.41zm305.713-1.62a2.324 2.324 0 11-3.149 3.42 2.324 2.324 0 013.149-3.42zm-200.074.135a2.324 2.324 0 11-3.42 3.149 2.324 2.324 0 013.42-3.149zm87.28-2.591a2.022 2.022 0 11-.508 4.012 2.022 2.022 0 01.507-4.012zm24.49 1.752a2.022 2.022 0 11-4.013.507 2.022 2.022 0 014.012-.507zm121.504-6a2.476 2.476 0 11-2.677 4.164 2.476 2.476 0 012.677-4.165zm-266.338.743a2.476 2.476 0 11-4.165 2.677 2.476 2.476 0 014.165-2.677zm206.964.134a2.173 2.173 0 11-3.298 2.831 2.173 2.173 0 013.298-2.83zm-147.946-.233a2.173 2.173 0 11-2.831 3.298 2.173 2.173 0 012.83-3.298zm107.713-.79a2.022 2.022 0 11-3.76 1.49 2.022 2.022 0 013.76-1.49zm-67.92-1.135a2.022 2.022 0 11-1.489 3.76 2.022 2.022 0 011.489-3.76zm-151.713-1.59a2.778 2.778 0 11-5.088 2.231 2.778 2.778 0 015.088-2.232zm372.38-1.429a2.778 2.778 0 11-2.23 5.088 2.778 2.778 0 012.23-5.088zm-22.75-4.573a2.627 2.627 0 11-2.237 4.754 2.627 2.627 0 012.237-4.754zM79.69 316.661a2.627 2.627 0 11-4.753 2.237 2.627 2.627 0 014.753-2.237zm276.953-3.652a2.325 2.325 0 11-2.543 3.892 2.325 2.325 0 012.543-3.892zm-227.14.674a2.325 2.325 0 11-3.892 2.543 2.325 2.325 0 013.892-2.543zm166.915-.951a2.022 2.022 0 11-3.272 2.377 2.022 2.022 0 013.272-2.377zm-107.082-.448a2.022 2.022 0 11-2.377 3.272 2.022 2.022 0 012.377-3.272zm268.024-2.301a2.93 2.93 0 11-1.865 5.553 2.93 2.93 0 011.865-5.553zm-428.082 1.843a2.93 2.93 0 11-5.553 1.866 2.93 2.93 0 015.553-1.866zm303.18-3.605a2.173 2.173 0 11-2.661 3.437 2.173 2.173 0 012.66-3.437zm-178.938.388a2.173 2.173 0 11-3.437 2.66 2.173 2.173 0 013.437-2.66zm100.987-1.278a1.87 1.87 0 11-3.69.616 1.87 1.87 0 013.69-.616zm-23.932-1.537a1.87 1.87 0 11-.616 3.69 1.87 1.87 0 01.616-3.69zM99.32 305.45a2.476 2.476 0 11-4.503 2.057 2.476 2.476 0 014.503-2.057zm287.567-1.223a2.476 2.476 0 11-2.057 4.503 2.476 2.476 0 012.057-4.503zm50.195-5.418a2.778 2.778 0 11-1.659 5.303 2.778 2.778 0 011.659-5.303zm-387.756 1.822a2.778 2.778 0 11-5.303 1.659 2.778 2.778 0 015.303-1.659zm226.162-1.154a1.87 1.87 0 11-3.29 1.78 1.87 1.87 0 013.29-1.78zm-65.513-.755a1.87 1.87 0 11-1.781 3.29 1.87 1.87 0 011.78-3.29zm-36.846-1.656a2.022 2.022 0 11-3.116 2.578 2.022 2.022 0 013.116-2.578zm139.516-.269a2.022 2.022 0 11-2.578 3.116 2.022 2.022 0 012.578-3.116zm102.168-2.699a2.627 2.627 0 11-1.623 4.997 2.627 2.627 0 011.623-4.997zm-343.39 1.687a2.627 2.627 0 11-4.996 1.623 2.627 2.627 0 014.997-1.623zm47.566-1.527a2.324 2.324 0 11-4.257 1.867 2.324 2.324 0 014.257-1.867zM367 293.063a2.324 2.324 0 11-1.867 4.257 2.324 2.324 0 011.867-4.257zm96.201-4.922a2.93 2.93 0 11-1.303 5.712 2.93 2.93 0 011.303-5.712zm-439.968 2.204a2.93 2.93 0 11-5.711 1.304 2.93 2.93 0 015.711-1.304zm320.986-1.826a2.173 2.173 0 11-1.914 3.902 2.173 2.173 0 011.914-3.902zm-202.603.994a2.173 2.173 0 11-3.902 1.914 2.173 2.173 0 013.902-1.914zm151.167-4.277a1.87 1.87 0 11-2.534 2.752 1.87 1.87 0 012.534-2.752zm-99.995.109a1.87 1.87 0 11-2.753 2.534 1.87 1.87 0 012.753-2.534zm-101.26-.785a2.476 2.476 0 11-4.75 1.395 2.476 2.476 0 014.75-1.395zm302.943-1.678a2.476 2.476 0 11-1.395 4.751 2.476 2.476 0 011.395-4.75zm-163.45.02a1.72 1.72 0 11-.823 3.34 1.72 1.72 0 01.823-3.34zm22.966 1.258a1.72 1.72 0 11-3.34.823 1.72 1.72 0 013.34-.823zm-210.106-5.03a2.778 2.778 0 11-5.453 1.065 2.778 2.778 0 015.453-1.065zm398.424-2.194a2.778 2.778 0 11-1.065 5.453 2.778 2.778 0 011.065-5.453zm-117.905.392a2.022 2.022 0 11-1.722 3.659 2.022 2.022 0 011.722-3.66zm-163.183.968a2.022 2.022 0 11-3.659 1.722 2.022 2.022 0 013.66-1.722zm-95.377-4.259a2.627 2.627 0 11-5.16.985 2.627 2.627 0 015.16-.985zm354.32-2.088a2.627 2.627 0 11-.985 5.16 2.627 2.627 0 01.985-5.16zm-46.226-.265a2.325 2.325 0 11-1.142 4.506 2.325 2.325 0 011.142-4.506zm-262.117 1.682a2.324 2.324 0 11-4.506 1.142 2.324 2.324 0 014.506-1.142zm161.01-.023a1.72 1.72 0 11-2.575 2.28 1.72 1.72 0 012.574-2.28zm-60.299-.147a1.72 1.72 0 11-2.28 2.575 1.72 1.72 0 012.28-2.575zm254.312-7.37a2.93 2.93 0 11-.729 5.813 2.93 2.93 0 01.728-5.812zM19.358 268.37a2.93 2.93 0 11-5.813.728 2.93 2.93 0 015.813-.728zm114.443.042a2.173 2.173 0 11-4.207 1.09 2.173 2.173 0 014.207-1.09zm217.974-1.559a2.173 2.173 0 11-1.09 4.208 2.173 2.173 0 011.09-4.208zm-47.259-.702a1.87 1.87 0 11-1.503 3.427 1.87 1.87 0 011.503-3.427zm-123.64.962a1.87 1.87 0 11-3.427 1.503 1.87 1.87 0 013.426-1.503zm-94.087-4.338a2.476 2.476 0 11-4.9.705 2.476 2.476 0 014.9-.705zm312.151-2.098a2.476 2.476 0 11-.705 4.901 2.476 2.476 0 01.705-4.9zm-157.476.968a1.569 1.569 0 110 3.137 1.569 1.569 0 010-3.137zm203.457-7.032l.17.009a2.778 2.778 0 11-.337-.008h.167zM40.836 257.16a2.778 2.778 0 11-5.537.459 2.778 2.778 0 015.537-.459zm290.11-1.614a2.022 2.022 0 11-.758 3.973 2.022 2.022 0 01.757-3.973zm-176.597 1.608a2.022 2.022 0 11-3.973.758 2.022 2.022 0 013.973-.758zm46.317-1.145a1.72 1.72 0 11-3.216 1.22 1.72 1.72 0 013.216-1.22zm83.815-.999a1.72 1.72 0 11-1.22 3.216 1.72 1.72 0 011.22-3.216zm-221.454-3.248a2.627 2.627 0 11-5.243.33 2.627 2.627 0 015.243-.33zm359.496-2.462l.166.005a2.627 2.627 0 11-.33 0l.164-.005zm-45.427.147l.159.008a2.324 2.324 0 11-.316-.005l.157-.003zm-268.913 2.132a2.324 2.324 0 11-4.633.384 2.324 2.324 0 014.633-.383zm115-.49a1.568 1.568 0 11-2.716 1.57 1.568 1.568 0 012.717-1.57zm38.705-.573a1.568 1.568 0 11-1.569 2.716 1.568 1.568 0 011.569-2.716zm206.35-7.255a2.93 2.93 0 11-.145 5.856 2.93 2.93 0 01.146-5.856zm-450.546 2.855a2.93 2.93 0 11-5.857.146 2.93 2.93 0 015.857-.146zm112.703.052a2.173 2.173 0 11-4.341.22 2.173 2.173 0 014.34-.22zm224.265-2.063l.155.003a2.173 2.173 0 11-.308.006l.153-.009zm-45.388.178l.146.006a1.87 1.87 0 11-.29-.002l.144-.004zm-133.742 1.716a1.87 1.87 0 11-3.73.309 1.87 1.87 0 013.73-.31zm-92.806-7.94a2.476 2.476 0 110 4.952 2.476 2.476 0 010-4.951zm317.48 0a2.476 2.476 0 110 4.952 2.476 2.476 0 010-4.951zm-158.74 1.06a1.417 1.417 0 110 2.834 1.417 1.417 0 010-2.835zm206.795-4.284a2.778 2.778 0 11-5.554.153 2.778 2.778 0 015.554-.153zm-410.735-2.7a2.778 2.778 0 11-.153 5.554 2.778 2.778 0 01.153-5.554zm113.387.683l.15.004a2.022 2.022 0 11-.299.003l.149-.007zm183.101 1.895a2.022 2.022 0 11-4.036.254 2.022 2.022 0 014.036-.254zM196.5 233.35l.148.011a1.72 1.72 0 11-.294-.01l.146-.001zm91.695 1.511a1.72 1.72 0 11-3.414.415 1.72 1.72 0 013.414-.415zm-227.79-8.344l.165.006a2.627 2.627 0 11-.33 0l.164-.006zm364.74 2.462a2.627 2.627 0 11-5.243.33 2.627 2.627 0 015.244-.33zM105.9 226.975l.16.008a2.324 2.324 0 11-.317-.005l.157-.003zm273.48 2.132a2.324 2.324 0 11-4.633.384 2.324 2.324 0 014.633-.384zm-116.917-.694a1.568 1.568 0 11-2.717 1.568 1.568 1.568 0 012.717-1.568zm-39.852-.575a1.568 1.568 0 11-1.569 2.717 1.568 1.568 0 011.569-2.716zm247.913-4.468a2.93 2.93 0 11-5.842.438 2.93 2.93 0 015.842-.438zm-455.15-2.71l.173.008a2.93 2.93 0 11-.342-.006l.17-.002zm160.601 1.361a1.87 1.87 0 11-.918 3.627 1.87 1.87 0 01.918-3.627zm133.254 1.354a1.87 1.87 0 11-3.627.92 1.87 1.87 0 013.627-.92zm46.463-.338a2.173 2.173 0 11-4.296.658 2.173 2.173 0 014.296-.658zm-225.976-1.82a2.173 2.173 0 11-.658 4.297 2.173 2.173 0 01.658-4.297zm-45.023-5.723a2.475 2.475 0 11-.704 4.9 2.475 2.475 0 01.704-4.9zm316.347 2.098a2.476 2.476 0 11-4.9.705 2.476 2.476 0 014.9-.705zm-159.574-1.302a1.569 1.569 0 110 3.137 1.569 1.569 0 010-3.137zm-36.35-2.934a1.72 1.72 0 11-1.953 2.83 1.72 1.72 0 011.954-2.83zm75.09.438a1.72 1.72 0 11-2.83 1.954 1.72 1.72 0 012.83-1.954zm166.171-1.673a2.778 2.778 0 11-5.503.763 2.778 2.778 0 015.503-.763zm-406.69-2.37a2.778 2.778 0 11-.763 5.503 2.778 2.778 0 01.763-5.503zm116.134.83a2.022 2.022 0 11-1.25 3.847 2.022 2.022 0 011.25-3.846zm173.837 1.3a2.022 2.022 0 11-3.846 1.249 2.022 2.022 0 013.846-1.25zm45.96-5.318a2.324 2.324 0 11-4.506 1.141 2.324 2.324 0 014.506-1.14zm-265.481-1.682a2.324 2.324 0 11-1.141 4.506 2.324 2.324 0 011.14-4.506zm-46.383-.92a2.627 2.627 0 11-.985 5.16 2.627 2.627 0 01.985-5.16zm358.496 2.088a2.627 2.627 0 11-5.16.984 2.627 2.627 0 015.16-.984zm-122.264-3.746a1.87 1.87 0 11-3.133 2.046 1.87 1.87 0 013.133-2.046zm-114.45-.543a1.87 1.87 0 11-2.047 3.132 1.87 1.87 0 012.046-3.132zm282.14-1.112a2.93 2.93 0 11-5.769 1.017 2.93 2.93 0 015.77-1.017zm-449.028-2.376a2.93 2.93 0 11-1.018 5.77 2.93 2.93 0 011.018-5.77zm331.183 2.13a2.173 2.173 0 11-4.075 1.509 2.173 2.173 0 014.075-1.51zm-213.938-1.284a2.173 2.173 0 11-1.51 4.076 2.173 2.173 0 011.51-4.076zm127.449-.264a1.72 1.72 0 11-1.598 3.045 1.72 1.72 0 011.598-3.046zm-41.43.723a1.72 1.72 0 11-3.047 1.598 1.72 1.72 0 013.046-1.598zm-132.059-6.14a2.476 2.476 0 11-1.395 4.751 2.476 2.476 0 011.395-4.75zm306.298 1.679a2.476 2.476 0 11-4.75 1.395 2.476 2.476 0 014.75-1.395zm-154.685-1.655a1.72 1.72 0 110 3.44 1.72 1.72 0 010-3.44zm-75.505-3.237a2.022 2.022 0 11-2.167 3.415 2.022 2.022 0 012.167-3.415zm153.8.624a2.022 2.022 0 11-3.415 2.167 2.022 2.022 0 013.415-2.167zm122.247-1.097a2.778 2.778 0 11-5.386 1.364 2.778 2.778 0 015.386-1.364zm-397.71-2.01a2.778 2.778 0 11-1.363 5.385 2.778 2.778 0 011.363-5.386zm240.102-2.368a1.87 1.87 0 11-2.298 2.952 1.87 1.87 0 012.298-2.952zm-83.244.327a1.87 1.87 0 11-2.953 2.298 1.87 1.87 0 012.953-2.298zm167.041-.754a2.325 2.325 0 11-4.258 1.867 2.325 2.325 0 014.258-1.867zm-250.4-1.195a2.324 2.324 0 11-1.868 4.257 2.324 2.324 0 011.867-4.257zm-48.058-1.775a2.627 2.627 0 11-1.623 4.996 2.627 2.627 0 011.623-4.996zm346.763 1.687a2.627 2.627 0 11-4.997 1.623 2.627 2.627 0 014.997-1.623zm-76.968-4.254a2.173 2.173 0 11-3.689 2.299 2.173 2.173 0 013.69-2.3zm-193.143-.695a2.173 2.173 0 11-2.3 3.689 2.173 2.173 0 012.3-3.689zm316.19-.39a2.929 2.929 0 11-5.639 1.586 2.929 2.929 0 015.64-1.587zm-438.618-2.027a2.93 2.93 0 11-1.587 5.64 2.93 2.93 0 011.587-5.64zm197.182-.714a1.87 1.87 0 11-3.539 1.215 1.87 1.87 0 013.54-1.215zm43.018-1.162a1.87 1.87 0 11-1.215 3.539 1.87 1.87 0 011.215-3.54zm123.95-.857a2.476 2.476 0 11-4.504 2.057 2.476 2.476 0 014.504-2.057zM98.098 172.34a2.476 2.476 0 11-2.057 4.504 2.476 2.476 0 012.057-4.504zm206.845.598a2.022 2.022 0 11-2.768 2.948 2.022 2.022 0 012.768-2.948zm-124.099.09a2.022 2.022 0 11-2.948 2.768 2.022 2.022 0 012.948-2.768zm60.62-2.395a1.87 1.87 0 110 3.742 1.87 1.87 0 010-3.742zm-190.162-4.262a2.778 2.778 0 11-1.948 5.203 2.778 2.778 0 011.948-5.203zm383.9 1.627a2.778 2.778 0 11-5.202 1.948 2.778 2.778 0 015.203-1.948zm-306.374-3.828a2.324 2.324 0 11-2.542 3.892 2.324 2.324 0 012.542-3.892zm228.49.675a2.324 2.324 0 11-3.892 2.542 2.324 2.324 0 013.892-2.542zm-278.886-3.93a2.627 2.627 0 11-2.237 4.753 2.627 2.627 0 012.237-4.753zm329.56 1.258a2.627 2.627 0 11-4.753 2.237 2.627 2.627 0 014.754-2.237zm-247.211-1.332a2.173 2.173 0 11-2.995 3.15 2.173 2.173 0 012.995-3.15zm164.44.078a2.173 2.173 0 11-3.15 2.995 2.173 2.173 0 013.15-2.995zm-125.684-.846a2.022 2.022 0 11-3.544 1.948 2.022 2.022 0 013.544-1.948zm86.6-.798a2.022 2.022 0 11-1.947 3.544 2.022 2.022 0 011.948-3.544zm169.15-2.659a2.93 2.93 0 11-5.453 2.14 2.93 2.93 0 015.453-2.14zM31.439 154.96a2.93 2.93 0 11-2.14 5.453 2.93 2.93 0 012.14-5.453zm77.823-2.329a2.476 2.476 0 11-2.677 4.166 2.476 2.476 0 012.677-4.166zm267.826.745a2.476 2.476 0 11-4.165 2.676 2.476 2.476 0 014.165-2.676zm-112.562-2.658a2.022 2.022 0 11-1.006 3.917 2.022 2.022 0 011.006-3.917zm-43.661 1.456a2.022 2.022 0 11-3.917 1.005 2.022 2.022 0 013.917-1.005zm20.6-4.37a2.022 2.022 0 110 4.045 2.022 2.022 0 010-4.044zm184.582-.678a2.778 2.778 0 11-4.957 2.508 2.778 2.778 0 014.957-2.508zM60.615 145.9a2.778 2.778 0 11-2.509 4.958 2.778 2.778 0 012.509-4.958zm282.664.907a2.324 2.324 0 11-3.42 3.149 2.324 2.324 0 013.42-3.149zm-200.345-.136a2.324 2.324 0 11-3.149 3.42 2.324 2.324 0 013.149-3.42zm35.54-.441a2.173 2.173 0 11-3.567 2.483 2.173 2.173 0 013.568-2.483zm129.006-.542a2.173 2.173 0 11-2.483 3.567 2.173 2.173 0 012.483-3.567zm89.378-3.77a2.627 2.627 0 11-4.435 2.816 2.627 2.627 0 014.435-2.815zm-307.162-.81a2.627 2.627 0 11-2.815 4.436 2.627 2.627 0 012.815-4.436zm356.173-5.365a2.93 2.93 0 11-5.213 2.672 2.93 2.93 0 015.213-2.672zm-404.866-1.27a2.93 2.93 0 11-2.673 5.212 2.93 2.93 0 012.673-5.213zm82.115.238a2.476 2.476 0 11-3.242 3.742 2.476 2.476 0 013.242-3.742zm240.185.25a2.476 2.476 0 11-3.742 3.242 2.476 2.476 0 013.742-3.242zm-164.556.52a2.173 2.173 0 11-3.994 1.714 2.173 2.173 0 013.994-1.714zm88.289-1.14a2.173 2.173 0 11-1.714 3.994 2.173 2.173 0 011.714-3.995zm-127.309-2.607a2.324 2.324 0 11-3.668 2.855 2.324 2.324 0 013.668-2.855zm166.737-.407a2.324 2.324 0 11-2.856 3.669 2.324 2.324 0 012.856-3.669zm88.187-3.941a2.778 2.778 0 11-4.65 3.039 2.778 2.778 0 014.65-3.04zm-342.528-.806a2.778 2.778 0 11-3.039 4.65 2.778 2.778 0 013.039-4.65zm192.603.76a2.173 2.173 0 11-.874 4.258 2.173 2.173 0 01.874-4.257zm-43.957 1.692a2.173 2.173 0 11-4.258.875 2.173 2.173 0 014.258-.875zm20.696-4.057a2.173 2.173 0 110 4.347 2.173 2.173 0 010-4.347zm-138.11-2.104a2.627 2.627 0 11-3.35 4.048 2.627 2.627 0 013.35-4.048zm279.918.35a2.627 2.627 0 11-4.048 3.348 2.627 2.627 0 014.048-3.349zm-75.943-4.395a2.325 2.325 0 11-2.213 4.089 2.325 2.325 0 012.213-4.089zm-128.58.938a2.324 2.324 0 11-4.089 2.213 2.324 2.324 0 014.089-2.213zm-39.367-.819a2.476 2.476 0 11-3.742 3.243 2.476 2.476 0 013.742-3.242zm207.655-.25a2.476 2.476 0 11-3.242 3.743 2.476 2.476 0 013.242-3.742zm87.381-2.784a2.93 2.93 0 11-4.921 3.178 2.93 2.93 0 014.921-3.178zm-381.86-.871a2.93 2.93 0 11-3.178 4.92 2.93 2.93 0 013.179-4.92zm146.924-3.951a2.324 2.324 0 11-4.397 1.51 2.324 2.324 0 014.397-1.51zm86.916-1.444a2.324 2.324 0 11-1.51 4.397 2.324 2.324 0 011.51-4.397zm-200.712-1.003a2.778 2.778 0 11-3.532 4.29 2.778 2.778 0 013.532-4.29zm315.466.379a2.778 2.778 0 11-4.289 3.532 2.778 2.778 0 014.289-3.532zm-33.702-2.649a2.627 2.627 0 11-3.595 3.83 2.627 2.627 0 013.595-3.83zm-248.26.117a2.627 2.627 0 11-3.83 3.596 2.627 2.627 0 013.83-3.596zm209.433-1.578a2.476 2.476 0 11-2.676 4.165 2.476 2.476 0 012.676-4.165zm-170.898.744a2.476 2.476 0 11-4.166 2.677 2.476 2.476 0 014.166-2.677zm63.636.29a2.324 2.324 0 11-4.585.765 2.324 2.324 0 014.585-.765zm42.88-1.91a2.324 2.324 0 11-.765 4.585 2.324 2.324 0 01.765-4.585zm-22.777-1.888a2.324 2.324 0 110 4.649 2.324 2.324 0 010-4.649zm179.587-4.829a2.93 2.93 0 11-4.58 3.653 2.93 2.93 0 014.58-3.653zm-355.058-.463a2.93 2.93 0 11-3.653 4.58 2.93 2.93 0 013.653-4.58zm242.442-2.968a2.476 2.476 0 11-2.057 4.504 2.476 2.476 0 012.057-4.504zm-130.663 1.224a2.476 2.476 0 11-4.503 2.057 2.476 2.476 0 014.503-2.057zm207.943-2.86a2.778 2.778 0 11-3.874 3.983 2.778 2.778 0 013.874-3.982zm-284.575.055a2.778 2.778 0 11-3.982 3.874 2.778 2.778 0 013.982-3.874zm248.502-.666a2.627 2.627 0 11-3.088 4.25 2.627 2.627 0 013.088-4.25zm-212.688.58a2.627 2.627 0 11-4.25 3.089 2.627 2.627 0 014.25-3.088zm62.163-4.693a2.476 2.476 0 11-4.751 1.395 2.476 2.476 0 014.75-1.395zm87.766-1.678a2.476 2.476 0 11-1.395 4.75 2.476 2.476 0 011.395-4.75zm-22.476-4.89a2.476 2.476 0 11-.705 4.901 2.476 2.476 0 01.705-4.9zm-43.084 2.099a2.476 2.476 0 11-4.9.704 2.476 2.476 0 014.9-.704zm184.576-2.91a2.93 2.93 0 11-4.194 4.091 2.93 2.93 0 014.194-4.09zM81.172 80.1a2.93 2.93 0 11-4.09 4.193 2.93 2.93 0 014.09-4.194zm160.293-.78a2.476 2.476 0 110 4.952 2.476 2.476 0 010-4.951zm88.663-.063a2.627 2.627 0 11-2.53 4.604 2.627 2.627 0 012.53-4.604zm-173.76 1.037a2.627 2.627 0 11-4.604 2.53 2.627 2.627 0 014.604-2.53zm212.16-3.009a2.778 2.778 0 11-3.412 4.384 2.778 2.778 0 013.412-4.384zM118.3 77.77a2.778 2.778 0 11-4.384 3.413 2.778 2.778 0 014.384-3.413zm190.917-8.354a2.627 2.627 0 11-1.934 4.884 2.627 2.627 0 011.934-4.884zM177.123 70.89a2.627 2.627 0 11-4.885 1.934 2.627 2.627 0 014.885-1.934zm211.99-6.317a2.93 2.93 0 11-3.765 4.488 2.93 2.93 0 013.765-4.488zm-291.17.362a2.93 2.93 0 11-4.488 3.765 2.93 2.93 0 014.487-3.765zm251.855-.64a2.778 2.778 0 11-2.91 4.734 2.778 2.778 0 012.91-4.733zm-212.846.913a2.778 2.778 0 11-4.733 2.91 2.778 2.778 0 014.733-2.91zm61.94-1.045a2.627 2.627 0 11-5.089 1.307 2.627 2.627 0 015.09-1.307zm88.343-1.89a2.627 2.627 0 11-1.307 5.088 2.627 2.627 0 011.307-5.089zm-22.704-4.331a2.627 2.627 0 11-.658 5.212 2.627 2.627 0 01.658-5.212zm-43.198 2.276a2.627 2.627 0 11-5.212.659 2.627 2.627 0 015.212-.659zm20.132-3.728a2.627 2.627 0 110 5.254 2.627 2.627 0 010-5.254zm88.288-3.042a2.778 2.778 0 11-2.371 5.024 2.778 2.778 0 012.37-5.024zm-172.879 1.327a2.778 2.778 0 11-5.024 2.37 2.778 2.778 0 015.024-2.37zm213.986-4.028a2.93 2.93 0 11-3.3 4.84 2.93 2.93 0 013.3-4.84zm-254.72.77a2.93 2.93 0 11-4.84 3.3 2.93 2.93 0 014.84-3.3zm192.496-6.646a2.778 2.778 0 11-1.804 5.255 2.778 2.778 0 011.804-5.255zm-130.813 1.725a2.778 2.778 0 11-5.255 1.804 2.778 2.778 0 015.255-1.804zm173.497-7.79a2.93 2.93 0 11-2.802 5.145 2.93 2.93 0 012.802-5.145zM135.583 39.98a2.93 2.93 0 11-5.145 2.802 2.93 2.93 0 015.145-2.802zm151.12-1.309a2.778 2.778 0 11-1.214 5.422 2.778 2.778 0 011.215-5.422zm-87.16 2.103a2.778 2.778 0 11-5.42 1.215 2.778 2.778 0 015.42-1.215zm64.68-5.854a2.778 2.778 0 11-.612 5.522 2.778 2.778 0 01.611-5.522zm-42.45 2.455a2.778 2.778 0 11-5.522.611 2.778 2.778 0 015.523-.61zm19.692-3.711a2.778 2.778 0 110 5.556 2.778 2.778 0 010-5.556zm89.223-4.791a2.93 2.93 0 11-2.275 5.398 2.93 2.93 0 012.275-5.398zm-174.61 1.561a2.93 2.93 0 11-5.399 2.276 2.93 2.93 0 015.399-2.276zM309.17 21.04a2.93 2.93 0 11-1.727 5.598 2.93 2.93 0 011.727-5.598zm-131.749 1.935a2.929 2.929 0 11-5.597 1.727 2.929 2.929 0 015.597-1.727zM286.98 15.39a2.93 2.93 0 11-1.161 5.742 2.93 2.93 0 011.16-5.742zm-87.577 2.29a2.93 2.93 0 11-5.743 1.161 2.93 2.93 0 015.743-1.16zm64.933-5.703a2.93 2.93 0 11-.583 5.83 2.93 2.93 0 01.583-5.83zm-45.741 0a2.93 2.93 0 11.583 5.83 2.93 2.93 0 01-.583-5.83zm22.87-1.141a2.93 2.93 0 110 5.858 2.93 2.93 0 010-5.858z' fill='%23e5e5e5'/%3E%3C/svg%3E");
            background-repeat: no-repeat;
            background-position: center center;
            background-size: 482px;
            height: 100%;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        body {
            margin: 0;
            text-align: center;
            padding: 2rem;
        }
        h1 {
            margin: 0;
        }
    </style>
</head>

<body>
    <h1>{{.Message}}</h1>
</body>

</html>
`

// Data is the template data.
type Data struct {
	Message string
}

// Run the web server.
func Run(addr string, clientset kubernetes.Interface, refresh time.Duration, username, password string) error {
	pods := &PodList{}

	// Background task to refresh the list of Pods.
	go func(pods *PodList) {
		limiter := time.Tick(refresh)

		for {
			<-limiter

			list, err := getPodList(clientset)
			if err != nil {
				panic(err)
			}

			*pods = list
		}
	}(pods)

	tmpl, err := template.New("output").Parse(tmplSource)
	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		pod, exist := getPod(r.Host, pods)

		if !exist {
			w.WriteHeader(http.StatusServiceUnavailable)
			tmpl.Execute(w, Data{Message: "Environment cannot be found. Try rebuilding to reinstate."})
			return
		}

		if pod.Status == corev1.PodPending {
			w.WriteHeader(http.StatusServiceUnavailable)
			tmpl.Execute(w, Data{Message: "Environment is currently building. Check back shortly."})
			return
		}

		if pod.Status == corev1.PodFailed {
			w.WriteHeader(http.StatusServiceUnavailable)
			tmpl.Execute(w, Data{Message: "Environment is currently in a failed state."})
			return
		}

		if pod.Status == corev1.PodUnknown {
			w.WriteHeader(http.StatusServiceUnavailable)
			tmpl.Execute(w, Data{Message: "Environment status is unknown. Please try rebuilding to reinstate."})
			return
		}

		if pod.Status == corev1.PodSucceeded {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment has shutdown. Please rebuild to reinstate.")
			return
		}

		// @todo, Determine if we remove the hardcoded "http://".
		endpoint := fmt.Sprintf("http://%s", pod.IP)

		url, err := url.Parse(endpoint)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(url)

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)

		proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(addr, Auth(handler, username, password, "Restricted"))
}

func Auth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.UserAgent(), "Sajaribot") {
			return
		}

		if strings.HasPrefix(r.URL.Path, "/graphql") {
			return
		}

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}
