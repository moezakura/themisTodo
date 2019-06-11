<template>
	<div class="task-timer-simple-controller">
		<div class="task-start-button base-timer-button"><span class="start-button-text">Start Task</span><i
				class="fas fa-play start-button-icon"></i></div>
		<div class="timer-start-and-end-text">
			<div class="timer-start">00:00</div>
			<span class="timer-text-splitter">〜</span>
			<div class="timer-end">00:00</div>
		</div>
		<div class="timer-total-texts">
			<div class="timer-today-text"><span
					class="timer-today-text-label">today</span><span>:</span><span class="timer-today-text-content">24:00</span>
			</div>
			<div class="timer-total-text"><span
					class="timer-total-text-label">total</span><span>:</span><span class="timer-total-text-content">9999:59</span>
			</div>
		</div>
		<div class="task-timer-edit-button base-timer-button"><i class="fas fa-edit"></i></div>
		<div class="task-timer-history-button base-timer-button"><i class="fas fa-history"></i></div>
	</div>
</template>

<script>
    export default {
        name: "TaskTimerSimpleController"
    }
</script>

<style lang="scss" scoped>
	.task-timer-simple-controller {
		$height: 45px;
		display: flex;

		.base-timer-button {
			border: 1px solid rgba(67, 160, 71, 0.5);
			background-color: rgba(67, 160, 71, 0.05);
			transition: background-color ease .3s;

			&:hover {
				background-color: rgba(67, 160, 71, 0.2);
			}
		}

		.task-start-button {
			$button-height: 40px;
			display: flex;
			width: 180px;
			height: $button-height;
			line-height: $button-height;
			text-align: center;
			font-size: 14px;
			margin: #{($height - $button-height) / 2} 0;

			.start-button-text {
				margin: 0 5px 0 auto;
				vertical-align: middle;
			}

			.start-button-icon {
				line-height: $button-height;
				text-align: center;
				margin: 0 auto 0 5px;
				vertical-align: middle;
			}
		}

		.timer-start-and-end-text {
			margin: 0 5px 0 auto;
			display: flex;
			line-height: $height;
			text-align: center;

			.timer-text-splitter {
				width: 35px;
			}

			.timer-start, .timer-end {
				width: 80px;
				letter-spacing: 1px;
			}
		}

		.timer-total-texts {
			margin: 0 auto 0 5px;
			line-height: $height / 2;
			text-align: left;
			letter-spacing: 1px;
			width: 110px;

			.timer-today-text,
			.timer-total-text {
				width: 100%;
				display: flex;

				span {
					display: block;
				}

				.timer-today-text-label,
				.timer-total-text-label {
					width: 38%;
				}

				.timer-today-text-content,
				.timer-total-text-content {
					margin-left: 5px;
				}
			}
		}

		.task-timer-edit-button,
		.task-timer-history-button {
			$button-height: 40px;
			width: 60px;
			height: $button-height;
			line-height: $button-height;
			margin: #{($height - $button-height) / 2} 0;
			text-align: center;
			font-size: 16px;
		}

		.task-timer-history-button {
			margin-right: 0;
			border-left: 0;
		}
	}
</style>