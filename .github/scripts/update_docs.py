import re
import os
import time
import subprocess
import requests
from google import genai
from google.genai import types

def get_jira_ticket_from_branch():
    branch = subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"]).decode("utf-8").strip()
    match = re.match(r"[A-Z]+-\d+", branch)
    return match.group(0) if match else None

def fetch_jira_and_notion_url(ticket_id):
    url = f"{os.getenv('JIRA_BASE_URL')}/rest/api/3/issue/{ticket_id}"
    auth = (os.getenv('JIRA_USER_EMAIL'), os.getenv('JIRA_API_TOKEN'))
    res = requests.get(url, auth=auth).json()

    summary = res['fields']['summary']
    description = str(res['fields'].get('description', ''))

    notion_match = re.search(r"notion.so/[^\s]+", description)
    notion_url = notion_match.group(0) if notion_match else None

    return summary, description, notion_url

def fetch_notion_markdown(notion_url):
    if not notion_url:
        return "No FSD provided"
    page_id = notion_url.split("-")[-1].split("?")[0]
    headers = {
        "Authorization": f"Bearer {os.getenv('NOTION_API_TOKEN')}",
        "Notion-Version": "2022-06-28"
    }

    res = requests.get(f"api.notion.com/v1/blocks/{page_id}/children", headers=headers).json()
    text_content = []
    for block in res.get('results', []):
        if block.get('type') == 'paragraph':
            rich_text = block.get('paragraph', {}).get('rich_text') or []
            if rich_text:
                text_content.append(rich_text[0].get('plain_text', ''))
    return "\n".join(text_content)

def get_git_diff():
    return subprocess.check_output(["git", "diff", "origin/main...HEAD"]).decode("utf-8")

def update_documentation(context_data):
    client = genai.Client(api_key=os.getenv('GEMINI_API_KEY'))

    system_prompt = (
        "You are an automated technical writer inside a CI/CD pipeline.\n"
        "Your task is to review an engineering changes package (Git Diff, Jira Ticket, and FSD) "
        "and update the relevant Markdown files inside the /docs folder.\n"
        "If key docs do not exist yet, create files such as: docs/database_schema.md, docs/logic_flows.md, "
        "docs/apis_called.md, docs/service.md, and docs/architecture.md.\n"
        "Output ONLY valid markdown files wrapped in standard markdown blocks, prefixed with "
        "the filepath like: File: docs/architecture.md\n"
        "Do not delete existing unrelated documentation blocks.\n"
        "If the PR diff is the only source of truth, synthesize technical specifications from code changes."
    )

    user_prompt = f"""
    === JIRA SUMMARY ===
    {context_data['jira_summary']}

    === JIRA DESCRIPTION ===
    {context_data['jira_desc']}

    === NOTION FSD TEXT ===
    {context_data['fsd_text']}

    === GIT DIFF ===
    {context_data['git_diff']}
    """

    response = client.models.generate_content(
        model='gemini-2.5-pro',
        contents=user_prompt,
        config=types.GenerateContentConfig(
            system_instruction=system_prompt,
            temperature=0.2,
        )
    )

    return response.text

def write_files_to_disk(llm_output):
    files = re.split(r"File:\s*(docs/[a-zA-Z0-9_\-\.\/]+)", llm_output)

    for i in range(1, len(files), 2):
        path = files[i].strip()
        content_block = files[i+1].strip()

        clean_content = re.sub(r"^```markdown\n|^```\n|```$", "", content_block, flags=re.MULTILINE)

        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as f:
            f.write(clean_content)
        print(f"Updated: {path}")


def get_commit_list(base_sha, head_sha):
    """Return commits in base_sha..head_sha oldest-first, with parsed ticket IDs."""
    out = subprocess.check_output(
        ["git", "log", "--reverse", "--pretty=%H%x1f%s", f"{base_sha}..{head_sha}"]
    ).decode("utf-8").strip()
    commits = []
    if not out:
        return commits
    for line in out.splitlines():
        sha, subject = line.split("\x1f", 1)
        ticket_match = re.search(r"[A-Z]+-\d+", subject)
        commits.append({
            "sha": sha,
            "subject": subject,
            "ticket": ticket_match.group(0) if ticket_match else None,
        })
    return commits


def chunk(seq, size):
    for i in range(0, len(seq), size):
        yield seq[i:i + size]


def get_batch_diff(prev_sha, last_sha):
    """Diff for a batch, excluding docs/ (to avoid feedback loops) and noisy paths."""
    return subprocess.check_output([
        "git", "diff", prev_sha, last_sha,
        "--",
        ".",
        ":(exclude)docs/**",
        ":(exclude)**/*.lock",
        ":(exclude)**/dist/**",
    ]).decode("utf-8", errors="replace")


def read_current_docs():
    """Snapshot the current on-disk docs so each batch sees prior batches' edits."""
    docs = {}
    if not os.path.isdir("docs"):
        return docs
    for root, _, files in os.walk("docs"):
        for name in files:
            if name.endswith(".md"):
                path = os.path.join(root, name)
                try:
                    with open(path, "r") as f:
                        docs[path] = f.read()
                except OSError:
                    pass
    return docs


def build_batch_context(batch, prev_sha, last_sha):
    """Pack batch metadata into the payload shape update_documentation already expects."""
    tickets = sorted({c["ticket"] for c in batch if c["ticket"]})

    jira_blocks = []
    for ticket in tickets:
        try:
            summary, desc, notion_url = fetch_jira_and_notion_url(ticket)
            fsd = fetch_notion_markdown(notion_url)
            jira_blocks.append(
                f"### {ticket}: {summary}\n{desc}\n\nFSD:\n{fsd}"
            )
        except Exception as e:
            jira_blocks.append(f"### {ticket}: (failed to fetch context: {e})")

    commit_log = "\n".join(f"- {c['sha'][:8]} {c['subject']}" for c in batch)
    current_docs = read_current_docs()
    docs_snapshot = "\n\n".join(
        f"File: {path}\n```markdown\n{content}\n```"
        for path, content in current_docs.items()
    ) or "(no existing docs)"

    return {
        "jira_summary": (
            f"Batch of {len(batch)} commits ({prev_sha[:8]}..{last_sha[:8]}) "
            f"touching tickets: {', '.join(tickets) if tickets else 'none'}"
        ),
        "jira_desc": "\n\n".join(jira_blocks) or "No ticket references in this batch.",
        "fsd_text": (
            f"COMMITS IN THIS BATCH:\n{commit_log}\n\n"
            f"CURRENT DOCS STATE (update these, do not overwrite unrelated content):\n"
            f"{docs_snapshot}"
        ),
        "git_diff": get_batch_diff(prev_sha, last_sha),
    }


def run_release_pr_flow(base_sha, head_sha, batch_size):
    commits = get_commit_list(base_sha, head_sha)
    if not commits:
        print("No commits in range; nothing to document.")
        return

    total_batches = (len(commits) + batch_size - 1) // batch_size
    print(f"Release PR: {len(commits)} commits, {total_batches} batches of up to {batch_size}")

    prev_sha = base_sha
    for i, batch in enumerate(chunk(commits, batch_size), start=1):
        last_sha = batch[-1]["sha"]
        print(f"\n=== Batch {i}/{total_batches}: {prev_sha[:8]}..{last_sha[:8]} ({len(batch)} commits) ===")
        try:
            payload = build_batch_context(batch, prev_sha, last_sha)
            output = update_documentation(payload)
            write_files_to_disk(output)
        except Exception as e:
            print(f"Batch {i} failed, continuing: {e}")
        prev_sha = last_sha
        time.sleep(2)


def run_feature_pr_flow():
    ticket = get_jira_ticket_from_branch()
    if not ticket:
        print("No Jira ticket found in branch name. Skipping automation.")
        return

    summary, desc, notion_link = fetch_jira_and_notion_url(ticket)
    fsd = fetch_notion_markdown(notion_link)
    diff = get_git_diff()

    payload = {
        "jira_summary": summary,
        "jira_desc": desc,
        "fsd_text": fsd,
        "git_diff": diff,
    }
    output = update_documentation(payload)
    write_files_to_disk(output)


if __name__ == "__main__":
    base_ref = os.getenv("GITHUB_BASE_REF", "")
    head_ref = os.getenv("GITHUB_HEAD_REF", "")
    base_sha = os.getenv("PR_BASE_SHA")
    head_sha = os.getenv("PR_HEAD_SHA")
    batch_size = int(os.getenv("BATCH_SIZE", "5"))

    is_release_pr = base_ref == "stg" and base_sha and head_sha

    if is_release_pr:
        run_release_pr_flow(base_sha, head_sha, batch_size)
    else:
        run_feature_pr_flow()