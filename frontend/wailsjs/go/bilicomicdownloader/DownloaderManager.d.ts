// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {bilicomicdownloader} from '../models';

export function ClearDownloaders():Promise<void>;

export function DownloadList(arg1:Array<number>):Promise<void>;

export function GetBookInfo():Promise<bilicomicdownloader.BookInfo>;

export function GetChapter():Promise<Array<bilicomicdownloader.Volume>>;

export function GetDownloader(arg1:string):Promise<void>;

export function GetDownloaders():Promise<Array<bilicomicdownloader.DownloaderSingle>>;

export function MessageSend(arg1:string):Promise<void>;

export function ProcessSend():Promise<void>;
