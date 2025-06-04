import { useNotify } from "./useNotification";

export function useRunCommand() {
  const notify = useNotify();

  async function runCommand<T>(params: {
    command: () => Promise<T>,
    onSuccess?: (data: T) => void,
    errMsg?: string,
    onError?: (err: any) => void
  }) {
    const {
      command,
      onSuccess,
      errMsg = "请求失败",
      onError
    } = params;

    try {
      const res = await command();
      onSuccess?.(res);
    } catch (err) {
      console.error(`${errMsg}：`, err);
      if (onError) {
        onError(err);
      } else {
        notify.error({ content: `${errMsg}：` + err });
      }
    }
  }

  return runCommand;
}
